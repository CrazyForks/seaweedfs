package s3api

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/private/protocol/xml/xmlutil"
	"github.com/seaweedfs/seaweedfs/weed/s3api/s3bucket"

	"github.com/seaweedfs/seaweedfs/weed/filer"
	"github.com/seaweedfs/seaweedfs/weed/s3api/s3_constants"
	"github.com/seaweedfs/seaweedfs/weed/storage/needle"

	"github.com/seaweedfs/seaweedfs/weed/s3api/s3err"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/seaweedfs/seaweedfs/weed/glog"
	"github.com/seaweedfs/seaweedfs/weed/pb/filer_pb"
	util_http "github.com/seaweedfs/seaweedfs/weed/util/http"
)

func (s3a *S3ApiServer) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {

	glog.V(3).Infof("ListBucketsHandler")

	var identity *Identity
	var s3Err s3err.ErrorCode
	if s3a.iam.isEnabled() {
		// Use authRequest instead of authUser for consistency with other endpoints
		// This ensures the same authentication flow and any fixes (like prefix handling) are applied
		identity, s3Err = s3a.iam.authRequest(r, s3_constants.ACTION_LIST)
		if s3Err != s3err.ErrNone {
			s3err.WriteErrorResponse(w, r, s3Err)
			return
		}
	}

	var response ListAllMyBucketsResult

	entries, _, err := s3a.list(s3a.option.BucketsPath, "", "", false, math.MaxInt32)

	if err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}

	identityId := r.Header.Get(s3_constants.AmzIdentityId)

	var listBuckets ListAllMyBucketsList
	for _, entry := range entries {
		if entry.IsDirectory {
			if identity != nil && !identity.canDo(s3_constants.ACTION_LIST, entry.Name, "") {
				continue
			}
			listBuckets.Bucket = append(listBuckets.Bucket, ListAllMyBucketsEntry{
				Name:         entry.Name,
				CreationDate: time.Unix(entry.Attributes.Crtime, 0).UTC(),
			})
		}
	}

	response = ListAllMyBucketsResult{
		Owner: CanonicalUser{
			ID:          identityId,
			DisplayName: identityId,
		},
		Buckets: listBuckets,
	}

	writeSuccessResponseXML(w, r, response)
}

func (s3a *S3ApiServer) PutBucketHandler(w http.ResponseWriter, r *http.Request) {

	// collect parameters
	bucket, _ := s3_constants.GetBucketAndObject(r)

	// validate the bucket name
	err := s3bucket.VerifyS3BucketName(bucket)
	if err != nil {
		glog.Errorf("put invalid bucket name: %v %v", bucket, err)
		s3err.WriteErrorResponse(w, r, s3err.ErrInvalidBucketName)
		return
	}

	// avoid duplicated buckets
	errCode := s3err.ErrNone
	if err := s3a.WithFilerClient(false, func(client filer_pb.SeaweedFilerClient) error {
		if resp, err := client.CollectionList(context.Background(), &filer_pb.CollectionListRequest{
			IncludeEcVolumes:     true,
			IncludeNormalVolumes: true,
		}); err != nil {
			glog.Errorf("list collection: %v", err)
			return fmt.Errorf("list collections: %w", err)
		} else {
			for _, c := range resp.Collections {
				if s3a.getCollectionName(bucket) == c.Name {
					errCode = s3err.ErrBucketAlreadyExists
					break
				}
			}
		}
		return nil
	}); err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}
	if exist, err := s3a.exists(s3a.option.BucketsPath, bucket, true); err == nil && exist {
		errCode = s3err.ErrBucketAlreadyExists
	}
	if errCode != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, errCode)
		return
	}

	fn := func(entry *filer_pb.Entry) {
		if identityId := r.Header.Get(s3_constants.AmzIdentityId); identityId != "" {
			if entry.Extended == nil {
				entry.Extended = make(map[string][]byte)
			}
			entry.Extended[s3_constants.AmzIdentityId] = []byte(identityId)
		}
	}

	// create the folder for bucket, but lazily create actual collection
	if err := s3a.mkdir(s3a.option.BucketsPath, bucket, fn); err != nil {
		glog.Errorf("PutBucketHandler mkdir: %v", err)
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}

	// Check for x-amz-bucket-object-lock-enabled header (S3 standard compliance)
	if objectLockHeaderValue := r.Header.Get(s3_constants.AmzBucketObjectLockEnabled); strings.EqualFold(objectLockHeaderValue, "true") {
		glog.V(3).Infof("PutBucketHandler: enabling Object Lock and Versioning for bucket %s due to x-amz-bucket-object-lock-enabled header", bucket)

		// Atomically update the configuration of the specified bucket. See the updateBucketConfig
		// function definition for detailed documentation on parameters and behavior.
		errCode := s3a.updateBucketConfig(bucket, func(bucketConfig *BucketConfig) error {
			// Enable versioning (required for Object Lock)
			bucketConfig.Versioning = s3_constants.VersioningEnabled

			// Create basic Object Lock configuration (enabled without default retention)
			objectLockConfig := &ObjectLockConfiguration{
				ObjectLockEnabled: s3_constants.ObjectLockEnabled,
			}

			// Set the cached Object Lock configuration
			bucketConfig.ObjectLockConfig = objectLockConfig

			return nil
		})

		if errCode != s3err.ErrNone {
			glog.Errorf("PutBucketHandler: failed to enable Object Lock for bucket %s: %v", bucket, errCode)
			s3err.WriteErrorResponse(w, r, errCode)
			return
		}
		glog.V(3).Infof("PutBucketHandler: enabled Object Lock and Versioning for bucket %s", bucket)
	}

	w.Header().Set("Location", "/"+bucket)
	writeSuccessResponseEmpty(w, r)
}

func (s3a *S3ApiServer) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {

	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("DeleteBucketHandler %s", bucket)

	if err := s3a.checkBucket(r, bucket); err != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, err)
		return
	}

	err := s3a.WithFilerClient(false, func(client filer_pb.SeaweedFilerClient) error {
		if !s3a.option.AllowDeleteBucketNotEmpty {
			entries, _, err := s3a.list(s3a.option.BucketsPath+"/"+bucket, "", "", false, 2)
			if err != nil {
				return fmt.Errorf("failed to list bucket %s: %v", bucket, err)
			}
			for _, entry := range entries {
				if entry.Name != s3_constants.MultipartUploadsFolder {
					return errors.New(s3err.GetAPIError(s3err.ErrBucketNotEmpty).Code)
				}
			}
		}

		// delete collection
		deleteCollectionRequest := &filer_pb.DeleteCollectionRequest{
			Collection: s3a.getCollectionName(bucket),
		}

		glog.V(1).Infof("delete collection: %v", deleteCollectionRequest)
		if _, err := client.DeleteCollection(context.Background(), deleteCollectionRequest); err != nil {
			return fmt.Errorf("delete collection %s: %v", bucket, err)
		}

		return nil
	})

	if err != nil {
		s3ErrorCode := s3err.ErrInternalError
		if err.Error() == s3err.GetAPIError(s3err.ErrBucketNotEmpty).Code {
			s3ErrorCode = s3err.ErrBucketNotEmpty
		}
		s3err.WriteErrorResponse(w, r, s3ErrorCode)
		return
	}

	err = s3a.rm(s3a.option.BucketsPath, bucket, false, true)

	if err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}

	s3err.WriteEmptyResponse(w, r, http.StatusNoContent)
}

func (s3a *S3ApiServer) HeadBucketHandler(w http.ResponseWriter, r *http.Request) {

	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("HeadBucketHandler %s", bucket)

	if entry, err := s3a.getEntry(s3a.option.BucketsPath, bucket); entry == nil || errors.Is(err, filer_pb.ErrNotFound) {
		s3err.WriteErrorResponse(w, r, s3err.ErrNoSuchBucket)
		return
	}

	writeSuccessResponseEmpty(w, r)
}

func (s3a *S3ApiServer) checkBucket(r *http.Request, bucket string) s3err.ErrorCode {
	entry, err := s3a.getEntry(s3a.option.BucketsPath, bucket)
	if entry == nil || errors.Is(err, filer_pb.ErrNotFound) {
		return s3err.ErrNoSuchBucket
	}

	//if iam is enabled, the access was already checked before
	if s3a.iam.isEnabled() {
		return s3err.ErrNone
	}
	if !s3a.hasAccess(r, entry) {
		return s3err.ErrAccessDenied
	}
	return s3err.ErrNone
}

func (s3a *S3ApiServer) hasAccess(r *http.Request, entry *filer_pb.Entry) bool {
	// Check if user is properly authenticated as admin through IAM system
	if s3a.isUserAdmin(r) {
		return true
	}

	if entry.Extended == nil {
		return true
	}

	identityId := r.Header.Get(s3_constants.AmzIdentityId)
	if id, ok := entry.Extended[s3_constants.AmzIdentityId]; ok {
		if identityId != string(id) {
			glog.V(3).Infof("hasAccess: %s != %s (entry.Extended = %v)", identityId, id, entry.Extended)
			return false
		}
	}
	return true
}

// isUserAdmin securely checks if the authenticated user is an admin
// This validates admin status through proper IAM authentication, not spoofable headers
func (s3a *S3ApiServer) isUserAdmin(r *http.Request) bool {
	// Use a minimal admin action to authenticate and check admin status
	adminAction := Action("Admin")
	identity, errCode := s3a.iam.authRequest(r, adminAction)
	if errCode != s3err.ErrNone {
		return false
	}

	// Check if the authenticated identity has admin privileges
	return identity != nil && identity.isAdmin()
}

// isBucketPublicRead checks if a bucket allows anonymous read access based on its cached ACL status
func (s3a *S3ApiServer) isBucketPublicRead(bucket string) bool {
	// Get bucket configuration which contains cached public-read status
	config, errCode := s3a.getBucketConfig(bucket)
	if errCode != s3err.ErrNone {
		return false
	}

	// Return the cached public-read status (no JSON parsing needed)
	return config.IsPublicRead
}

// isPublicReadGrants checks if the grants allow public read access
func isPublicReadGrants(grants []*s3.Grant) bool {
	for _, grant := range grants {
		if grant.Grantee != nil && grant.Grantee.URI != nil && grant.Permission != nil {
			// Check for AllUsers group with Read permission
			if *grant.Grantee.URI == s3_constants.GranteeGroupAllUsers &&
				(*grant.Permission == s3_constants.PermissionRead || *grant.Permission == s3_constants.PermissionFullControl) {
				return true
			}
		}
	}
	return false
}

// AuthWithPublicRead creates an auth wrapper that allows anonymous access for public-read buckets
func (s3a *S3ApiServer) AuthWithPublicRead(handler http.HandlerFunc, action Action) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket, _ := s3_constants.GetBucketAndObject(r)
		authType := getRequestAuthType(r)
		isAnonymous := authType == authTypeAnonymous

		if isAnonymous {
			isPublic := s3a.isBucketPublicRead(bucket)

			if isPublic {
				handler(w, r)
				return
			}
		}
		s3a.iam.Auth(handler, action)(w, r) // Fallback to normal IAM auth
	}
}

// GetBucketAclHandler Get Bucket ACL
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketAcl.html
func (s3a *S3ApiServer) GetBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	// collect parameters
	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("GetBucketAclHandler %s", bucket)

	if err := s3a.checkBucket(r, bucket); err != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, err)
		return
	}

	amzAccountId := r.Header.Get(s3_constants.AmzAccountId)
	amzDisplayName := s3a.iam.GetAccountNameById(amzAccountId)
	response := AccessControlPolicy{
		Owner: CanonicalUser{
			ID:          amzAccountId,
			DisplayName: amzDisplayName,
		},
	}
	response.AccessControlList.Grant = append(response.AccessControlList.Grant, Grant{
		Grantee: Grantee{
			ID:          amzAccountId,
			DisplayName: amzDisplayName,
			Type:        "CanonicalUser",
			XMLXSI:      "CanonicalUser",
			XMLNS:       "http://www.w3.org/2001/XMLSchema-instance"},
		Permission: s3.PermissionFullControl,
	})
	writeSuccessResponseXML(w, r, response)
}

// PutBucketAclHandler Put bucket ACL
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketAcl.html //
func (s3a *S3ApiServer) PutBucketAclHandler(w http.ResponseWriter, r *http.Request) {
	// collect parameters
	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("PutBucketAclHandler %s", bucket)

	if err := s3a.checkBucket(r, bucket); err != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, err)
		return
	}

	// Get account information for ACL processing
	amzAccountId := r.Header.Get(s3_constants.AmzAccountId)

	// Get bucket ownership settings (these would be used for ownership validation in a full implementation)
	bucketOwnership := ""         // Default/simplified for now - in a full implementation this would be retrieved from bucket config
	bucketOwnerId := amzAccountId // Simplified - bucket owner is current account

	// Use the existing ACL parsing logic to handle both canned ACLs and XML body
	grants, errCode := ExtractAcl(r, s3a.iam, bucketOwnership, bucketOwnerId, amzAccountId, amzAccountId)
	if errCode != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, errCode)
		return
	}

	// Store the bucket ACL in bucket metadata
	errCode = s3a.updateBucketConfig(bucket, func(config *BucketConfig) error {
		if len(grants) > 0 {
			grantsBytes, err := json.Marshal(grants)
			if err != nil {
				glog.Errorf("PutBucketAclHandler: failed to marshal grants: %v", err)
				return err
			}
			config.ACL = grantsBytes
			// Cache the public-read status to avoid JSON parsing on every request
			config.IsPublicRead = isPublicReadGrants(grants)
		} else {
			config.ACL = nil
			config.IsPublicRead = false
		}
		config.Owner = amzAccountId
		return nil
	})

	if errCode != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, errCode)
		return
	}

	glog.V(3).Infof("PutBucketAclHandler: Successfully stored ACL for bucket %s with %d grants", bucket, len(grants))

	writeSuccessResponseEmpty(w, r)
}

// GetBucketLifecycleConfigurationHandler Get Bucket Lifecycle configuration
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketLifecycleConfiguration.html
func (s3a *S3ApiServer) GetBucketLifecycleConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	// collect parameters
	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("GetBucketLifecycleConfigurationHandler %s", bucket)

	if err := s3a.checkBucket(r, bucket); err != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, err)
		return
	}
	fc, err := filer.ReadFilerConf(s3a.option.Filer, s3a.option.GrpcDialOption, nil)
	if err != nil {
		glog.Errorf("GetBucketLifecycleConfigurationHandler: %s", err)
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}
	ttls := fc.GetCollectionTtls(s3a.getCollectionName(bucket))
	if len(ttls) == 0 {
		s3err.WriteErrorResponse(w, r, s3err.ErrNoSuchLifecycleConfiguration)
		return
	}

	response := Lifecycle{}
	// Sort locationPrefixes to ensure consistent ordering of lifecycle rules
	var locationPrefixes []string
	for locationPrefix := range ttls {
		locationPrefixes = append(locationPrefixes, locationPrefix)
	}
	sort.Strings(locationPrefixes)

	for _, locationPrefix := range locationPrefixes {
		internalTtl := ttls[locationPrefix]
		ttl, _ := needle.ReadTTL(internalTtl)
		days := int(ttl.Minutes() / 60 / 24)
		if days == 0 {
			continue
		}
		prefix, found := strings.CutPrefix(locationPrefix, fmt.Sprintf("%s/%s/", s3a.option.BucketsPath, bucket))
		if !found {
			continue
		}
		response.Rules = append(response.Rules, Rule{
			ID:         prefix,
			Status:     Enabled,
			Prefix:     Prefix{val: prefix, set: true},
			Expiration: Expiration{Days: days, set: true},
		})
	}

	writeSuccessResponseXML(w, r, response)
}

// PutBucketLifecycleConfigurationHandler Put Bucket Lifecycle configuration
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketLifecycleConfiguration.html
func (s3a *S3ApiServer) PutBucketLifecycleConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	// collect parameters
	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("PutBucketLifecycleConfigurationHandler %s", bucket)

	if err := s3a.checkBucket(r, bucket); err != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, err)
		return
	}

	lifeCycleConfig := Lifecycle{}
	if err := xmlDecoder(r.Body, &lifeCycleConfig, r.ContentLength); err != nil {
		glog.Warningf("PutBucketLifecycleConfigurationHandler xml decode: %s", err)
		s3err.WriteErrorResponse(w, r, s3err.ErrMalformedXML)
		return
	}

	fc, err := filer.ReadFilerConf(s3a.option.Filer, s3a.option.GrpcDialOption, nil)
	if err != nil {
		glog.Errorf("PutBucketLifecycleConfigurationHandler read filer config: %s", err)
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}
	collectionName := s3a.getCollectionName(bucket)
	collectionTtls := fc.GetCollectionTtls(collectionName)
	changed := false

	for _, rule := range lifeCycleConfig.Rules {
		if rule.Status != Enabled {
			continue
		}
		var rulePrefix string
		switch {
		case rule.Filter.Prefix.set:
			rulePrefix = rule.Filter.Prefix.val
		case rule.Prefix.set:
			rulePrefix = rule.Prefix.val
		case !rule.Expiration.Date.IsZero() || rule.Transition.Days > 0 || !rule.Transition.Date.IsZero():
			s3err.WriteErrorResponse(w, r, s3err.ErrNotImplemented)
			return
		}

		if rule.Expiration.Days == 0 {
			continue
		}

		locConf := &filer_pb.FilerConf_PathConf{
			LocationPrefix: fmt.Sprintf("%s/%s/%s", s3a.option.BucketsPath, bucket, rulePrefix),
			Collection:     collectionName,
			Ttl:            fmt.Sprintf("%dd", rule.Expiration.Days),
		}
		if ttl, ok := collectionTtls[locConf.LocationPrefix]; ok && ttl == locConf.Ttl {
			continue
		}
		if err := fc.AddLocationConf(locConf); err != nil {
			glog.Errorf("PutBucketLifecycleConfigurationHandler add location config: %s", err)
			s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
			return
		}
		changed = true
	}

	if changed {
		var buf bytes.Buffer
		if err := fc.ToText(&buf); err != nil {
			glog.Errorf("PutBucketLifecycleConfigurationHandler save config to text: %s", err)
			s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		}
		if err := s3a.WithFilerClient(false, func(client filer_pb.SeaweedFilerClient) error {
			return filer.SaveInsideFiler(client, filer.DirectoryEtcSeaweedFS, filer.FilerConfName, buf.Bytes())
		}); err != nil {
			glog.Errorf("PutBucketLifecycleConfigurationHandler save config inside filer: %s", err)
			s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
			return
		}
	}

	writeSuccessResponseEmpty(w, r)
}

// DeleteBucketLifecycleHandler Delete Bucket Lifecycle
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketLifecycle.html
func (s3a *S3ApiServer) DeleteBucketLifecycleHandler(w http.ResponseWriter, r *http.Request) {
	// collect parameters
	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("DeleteBucketLifecycleHandler %s", bucket)

	if err := s3a.checkBucket(r, bucket); err != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, err)
		return
	}

	fc, err := filer.ReadFilerConf(s3a.option.Filer, s3a.option.GrpcDialOption, nil)
	if err != nil {
		glog.Errorf("DeleteBucketLifecycleHandler read filer config: %s", err)
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}
	collectionTtls := fc.GetCollectionTtls(s3a.getCollectionName(bucket))
	changed := false
	for prefix, ttl := range collectionTtls {
		bucketPrefix := fmt.Sprintf("%s/%s/", s3a.option.BucketsPath, bucket)
		if strings.HasPrefix(prefix, bucketPrefix) && strings.HasSuffix(ttl, "d") {
			pathConf, found := fc.GetLocationConf(prefix)
			if found {
				pathConf.Ttl = ""
				fc.SetLocationConf(pathConf)
			}
			changed = true
		}
	}

	if changed {
		var buf bytes.Buffer
		if err := fc.ToText(&buf); err != nil {
			glog.Errorf("DeleteBucketLifecycleHandler save config to text: %s", err)
			s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		}
		if err := s3a.WithFilerClient(false, func(client filer_pb.SeaweedFilerClient) error {
			return filer.SaveInsideFiler(client, filer.DirectoryEtcSeaweedFS, filer.FilerConfName, buf.Bytes())
		}); err != nil {
			glog.Errorf("DeleteBucketLifecycleHandler save config inside filer: %s", err)
			s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
			return
		}
	}

	s3err.WriteEmptyResponse(w, r, http.StatusNoContent)
}

// GetBucketLocationHandler Get bucket location
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketLocation.html
func (s3a *S3ApiServer) GetBucketLocationHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := s3_constants.GetBucketAndObject(r)

	if err := s3a.checkBucket(r, bucket); err != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, err)
		return
	}

	writeSuccessResponseXML(w, r, CreateBucketConfiguration{})
}

// GetBucketRequestPaymentHandler Get bucket location
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketRequestPayment.html
func (s3a *S3ApiServer) GetBucketRequestPaymentHandler(w http.ResponseWriter, r *http.Request) {
	writeSuccessResponseXML(w, r, RequestPaymentConfiguration{Payer: "BucketOwner"})
}

// PutBucketOwnershipControls https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketOwnershipControls.html
func (s3a *S3ApiServer) PutBucketOwnershipControls(w http.ResponseWriter, r *http.Request) {
	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("PutBucketOwnershipControls %s", bucket)

	errCode := s3a.checkAccessByOwnership(r, bucket)
	if errCode != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, errCode)
		return
	}

	if r.Body == nil || r.Body == http.NoBody {
		s3err.WriteErrorResponse(w, r, s3err.ErrInvalidRequest)
		return
	}

	var v s3.OwnershipControls
	defer util_http.CloseRequest(r)

	err := xmlutil.UnmarshalXML(&v, xml.NewDecoder(r.Body), "")
	if err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrInvalidRequest)
		return
	}

	if len(v.Rules) != 1 {
		s3err.WriteErrorResponse(w, r, s3err.ErrInvalidRequest)
		return
	}

	printOwnership := true
	ownership := *v.Rules[0].ObjectOwnership
	switch ownership {
	case s3_constants.OwnershipObjectWriter:
	case s3_constants.OwnershipBucketOwnerPreferred:
	case s3_constants.OwnershipBucketOwnerEnforced:
		printOwnership = false
	default:
		s3err.WriteErrorResponse(w, r, s3err.ErrInvalidRequest)
		return
	}

	// Check if ownership needs to be updated
	currentOwnership, errCode := s3a.getBucketOwnership(bucket)
	if errCode != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, errCode)
		return
	}

	if currentOwnership != ownership {
		errCode = s3a.setBucketOwnership(bucket, ownership)
		if errCode != s3err.ErrNone {
			s3err.WriteErrorResponse(w, r, errCode)
			return
		}
	}

	if printOwnership {
		result := &s3.PutBucketOwnershipControlsInput{
			OwnershipControls: &v,
		}
		s3err.WriteAwsXMLResponse(w, r, http.StatusOK, result)
	} else {
		writeSuccessResponseEmpty(w, r)
	}
}

// GetBucketOwnershipControls https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketOwnershipControls.html
func (s3a *S3ApiServer) GetBucketOwnershipControls(w http.ResponseWriter, r *http.Request) {
	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("GetBucketOwnershipControls %s", bucket)

	errCode := s3a.checkAccessByOwnership(r, bucket)
	if errCode != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, errCode)
		return
	}

	// Get ownership using new bucket config system
	ownership, errCode := s3a.getBucketOwnership(bucket)
	if errCode == s3err.ErrNoSuchBucket {
		s3err.WriteErrorResponse(w, r, s3err.ErrNoSuchBucket)
		return
	} else if errCode != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, s3err.OwnershipControlsNotFoundError)
		return
	}

	result := &s3.PutBucketOwnershipControlsInput{
		OwnershipControls: &s3.OwnershipControls{
			Rules: []*s3.OwnershipControlsRule{
				{
					ObjectOwnership: &ownership,
				},
			},
		},
	}

	s3err.WriteAwsXMLResponse(w, r, http.StatusOK, result)
}

// DeleteBucketOwnershipControls https://docs.aws.amazon.com/AmazonS3/latest/API/API_DeleteBucketOwnershipControls.html
func (s3a *S3ApiServer) DeleteBucketOwnershipControls(w http.ResponseWriter, r *http.Request) {
	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("PutBucketOwnershipControls %s", bucket)

	errCode := s3a.checkAccessByOwnership(r, bucket)
	if errCode != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, errCode)
		return
	}

	bucketEntry, err := s3a.getEntry(s3a.option.BucketsPath, bucket)
	if err != nil {
		if errors.Is(err, filer_pb.ErrNotFound) {
			s3err.WriteErrorResponse(w, r, s3err.ErrNoSuchBucket)
			return
		}
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}

	_, ok := bucketEntry.Extended[s3_constants.ExtOwnershipKey]
	if !ok {
		s3err.WriteErrorResponse(w, r, s3err.OwnershipControlsNotFoundError)
		return
	}

	delete(bucketEntry.Extended, s3_constants.ExtOwnershipKey)
	err = s3a.updateEntry(s3a.option.BucketsPath, bucketEntry)
	if err != nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrInternalError)
		return
	}

	emptyOwnershipControls := &s3.OwnershipControls{
		Rules: []*s3.OwnershipControlsRule{},
	}
	s3err.WriteAwsXMLResponse(w, r, http.StatusOK, emptyOwnershipControls)
}

// GetBucketVersioningHandler Get Bucket Versioning status
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_GetBucketVersioning.html
func (s3a *S3ApiServer) GetBucketVersioningHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("GetBucketVersioning %s", bucket)

	if err := s3a.checkBucket(r, bucket); err != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, err)
		return
	}

	// Get versioning status using new bucket config system
	versioningStatus, errCode := s3a.getBucketVersioningStatus(bucket)
	if errCode != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, errCode)
		return
	}

	// AWS S3 behavior: If versioning was never configured, don't return Status field
	var response *s3.PutBucketVersioningInput
	if versioningStatus == "" {
		// No versioning configuration - return empty response (no Status field)
		response = &s3.PutBucketVersioningInput{
			VersioningConfiguration: &s3.VersioningConfiguration{},
		}
	} else {
		// Versioning was explicitly configured - return the status
		response = &s3.PutBucketVersioningInput{
			VersioningConfiguration: &s3.VersioningConfiguration{
				Status: aws.String(versioningStatus),
			},
		}
	}
	s3err.WriteAwsXMLResponse(w, r, http.StatusOK, response)
}

// PutBucketVersioningHandler Put bucket Versioning
// https://docs.aws.amazon.com/AmazonS3/latest/API/API_PutBucketVersioning.html
func (s3a *S3ApiServer) PutBucketVersioningHandler(w http.ResponseWriter, r *http.Request) {
	bucket, _ := s3_constants.GetBucketAndObject(r)
	glog.V(3).Infof("PutBucketVersioning %s", bucket)

	if err := s3a.checkBucket(r, bucket); err != s3err.ErrNone {
		s3err.WriteErrorResponse(w, r, err)
		return
	}

	if r.Body == nil || r.Body == http.NoBody {
		s3err.WriteErrorResponse(w, r, s3err.ErrInvalidRequest)
		return
	}

	var versioningConfig s3.VersioningConfiguration
	defer util_http.CloseRequest(r)

	err := xmlutil.UnmarshalXML(&versioningConfig, xml.NewDecoder(r.Body), "")
	if err != nil {
		glog.Warningf("PutBucketVersioningHandler xml decode: %s", err)
		s3err.WriteErrorResponse(w, r, s3err.ErrMalformedXML)
		return
	}

	if versioningConfig.Status == nil {
		s3err.WriteErrorResponse(w, r, s3err.ErrInvalidRequest)
		return
	}

	status := *versioningConfig.Status
	if status != s3_constants.VersioningEnabled && status != s3_constants.VersioningSuspended {
		s3err.WriteErrorResponse(w, r, s3err.ErrInvalidRequest)
		return
	}

	// Check if trying to suspend versioning on a bucket with object lock enabled
	if status == s3_constants.VersioningSuspended {
		// Get bucket configuration to check for object lock
		bucketConfig, errCode := s3a.getBucketConfig(bucket)
		if errCode == s3err.ErrNone && bucketConfig.ObjectLockConfig != nil {
			// Object lock is enabled, cannot suspend versioning
			s3err.WriteErrorResponse(w, r, s3err.ErrInvalidBucketState)
			return
		}
	}

	// Update bucket versioning configuration using new bucket config system
	if errCode := s3a.setBucketVersioningStatus(bucket, status); errCode != s3err.ErrNone {
		glog.Errorf("PutBucketVersioningHandler save config: %d", errCode)
		s3err.WriteErrorResponse(w, r, errCode)
		return
	}

	writeSuccessResponseEmpty(w, r)
}
