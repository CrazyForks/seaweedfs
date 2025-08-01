package s3api

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/seaweedfs/seaweedfs/weed/pb/iam_pb"
	"github.com/seaweedfs/seaweedfs/weed/s3api/s3_constants"

	"github.com/seaweedfs/seaweedfs/weed/s3api/s3err"
	"github.com/stretchr/testify/assert"
)

// TestIsRequestPresignedSignatureV4 - Test validates the logic for presign signature version v4 detection.
func TestIsRequestPresignedSignatureV4(t *testing.T) {
	testCases := []struct {
		inputQueryKey   string
		inputQueryValue string
		expectedResult  bool
	}{
		// Test case - 1.
		// Test case with query key ""X-Amz-Credential" set.
		{"", "", false},
		// Test case - 2.
		{"X-Amz-Credential", "", true},
		// Test case - 3.
		{"X-Amz-Content-Sha256", "", false},
	}

	for i, testCase := range testCases {
		// creating an input HTTP request.
		// Only the query parameters are relevant for this particular test.
		inputReq, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
		if err != nil {
			t.Fatalf("Error initializing input HTTP request: %v", err)
		}
		q := inputReq.URL.Query()
		q.Add(testCase.inputQueryKey, testCase.inputQueryValue)
		inputReq.URL.RawQuery = q.Encode()

		actualResult := isRequestPresignedSignatureV4(inputReq)
		if testCase.expectedResult != actualResult {
			t.Errorf("Test %d: Expected the result to `%v`, but instead got `%v`", i+1, testCase.expectedResult, actualResult)
		}
	}
}

// Tests is requested authenticated function, tests replies for s3 errors.
func TestIsReqAuthenticated(t *testing.T) {
	iam := &IdentityAccessManagement{
		hashes:       make(map[string]*sync.Pool),
		hashCounters: make(map[string]*int32),
	}
	_ = iam.loadS3ApiConfiguration(&iam_pb.S3ApiConfiguration{
		Identities: []*iam_pb.Identity{
			{
				Name: "someone",
				Credentials: []*iam_pb.Credential{
					{
						AccessKey: "access_key_1",
						SecretKey: "secret_key_1",
					},
				},
				Actions: []string{"Read", "Write"},
			},
		},
	})

	// List of test cases for validating http request authentication.
	testCases := []struct {
		req     *http.Request
		s3Error s3err.ErrorCode
	}{
		// When request is unsigned, access denied is returned.
		{mustNewRequest(http.MethodGet, "http://127.0.0.1:9000", 0, nil, t), s3err.ErrAccessDenied},
		// When request is properly signed, error is none.
		{mustNewSignedRequest(http.MethodGet, "http://127.0.0.1:9000", 0, nil, t), s3err.ErrNone},
	}

	// Validates all testcases.
	for i, testCase := range testCases {
		if _, s3Error := iam.reqSignatureV4Verify(testCase.req); s3Error != testCase.s3Error {
			io.ReadAll(testCase.req.Body)
			t.Fatalf("Test %d: Unexpected S3 error: want %d - got %d", i, testCase.s3Error, s3Error)
		}
	}
}

func TestCheckaAnonymousRequestAuthType(t *testing.T) {
	iam := &IdentityAccessManagement{
		hashes:       make(map[string]*sync.Pool),
		hashCounters: make(map[string]*int32),
	}
	_ = iam.loadS3ApiConfiguration(&iam_pb.S3ApiConfiguration{
		Identities: []*iam_pb.Identity{
			{
				Name:    "anonymous",
				Actions: []string{s3_constants.ACTION_READ},
			},
		},
	})
	testCases := []struct {
		Request *http.Request
		ErrCode s3err.ErrorCode
		Action  Action
	}{
		{Request: mustNewRequest(http.MethodGet, "http://127.0.0.1:9000/bucket", 0, nil, t), ErrCode: s3err.ErrNone, Action: s3_constants.ACTION_READ},
		{Request: mustNewRequest(http.MethodPut, "http://127.0.0.1:9000/bucket", 0, nil, t), ErrCode: s3err.ErrAccessDenied, Action: s3_constants.ACTION_WRITE},
	}
	for i, testCase := range testCases {
		_, s3Error := iam.authRequest(testCase.Request, testCase.Action)
		if s3Error != testCase.ErrCode {
			t.Errorf("Test %d: Unexpected s3error returned wanted %d, got %d", i, testCase.ErrCode, s3Error)
		}
		if testCase.Request.Header.Get(s3_constants.AmzAuthType) != "Anonymous" {
			t.Errorf("Test %d: Unexpected AuthType returned wanted %s, got %s", i, "Anonymous", testCase.Request.Header.Get(s3_constants.AmzAuthType))
		}
	}

}

func TestCheckAdminRequestAuthType(t *testing.T) {
	iam := &IdentityAccessManagement{
		hashes:       make(map[string]*sync.Pool),
		hashCounters: make(map[string]*int32),
	}
	_ = iam.loadS3ApiConfiguration(&iam_pb.S3ApiConfiguration{
		Identities: []*iam_pb.Identity{
			{
				Name: "someone",
				Credentials: []*iam_pb.Credential{
					{
						AccessKey: "access_key_1",
						SecretKey: "secret_key_1",
					},
				},
				Actions: []string{"Admin", "Read", "Write"},
			},
		},
	})
	testCases := []struct {
		Request *http.Request
		ErrCode s3err.ErrorCode
	}{
		{Request: mustNewRequest(http.MethodGet, "http://127.0.0.1:9000", 0, nil, t), ErrCode: s3err.ErrAccessDenied},
		{Request: mustNewSignedRequest(http.MethodGet, "http://127.0.0.1:9000", 0, nil, t), ErrCode: s3err.ErrNone},
		{Request: mustNewPresignedRequest(iam, http.MethodGet, "http://127.0.0.1:9000", 0, nil, t), ErrCode: s3err.ErrNone},
	}
	for i, testCase := range testCases {
		if _, s3Error := iam.reqSignatureV4Verify(testCase.Request); s3Error != testCase.ErrCode {
			t.Errorf("Test %d: Unexpected s3error returned wanted %d, got %d", i, testCase.ErrCode, s3Error)
		}
	}
}

func BenchmarkGetSignature(b *testing.B) {
	t := time.Now()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		signingKey := getSigningKey("secret-key", t.Format(yyyymmdd), "us-east-1", "s3")
		getSignature(signingKey, "random data")
	}
}

// Provides a fully populated http request instance, fails otherwise.
func mustNewRequest(method string, urlStr string, contentLength int64, body io.ReadSeeker, t *testing.T) *http.Request {
	req, err := newTestRequest(method, urlStr, contentLength, body)
	if err != nil {
		t.Fatalf("Unable to initialize new http request %s", err)
	}
	return req
}

// This is similar to mustNewRequest but additionally the request
// is signed with AWS Signature V4, fails if not able to do so.
func mustNewSignedRequest(method string, urlStr string, contentLength int64, body io.ReadSeeker, t *testing.T) *http.Request {
	req := mustNewRequest(method, urlStr, contentLength, body, t)
	cred := &Credential{"access_key_1", "secret_key_1"}
	if err := signRequestV4(req, cred.AccessKey, cred.SecretKey); err != nil {
		t.Fatalf("Unable to initialized new signed http request %s", err)
	}
	return req
}

// This is similar to mustNewRequest but additionally the request
// is presigned with AWS Signature V4, fails if not able to do so.
func mustNewPresignedRequest(iam *IdentityAccessManagement, method string, urlStr string, contentLength int64, body io.ReadSeeker, t *testing.T) *http.Request {
	req := mustNewRequest(method, urlStr, contentLength, body, t)
	cred := &Credential{"access_key_1", "secret_key_1"}
	if err := preSignV4(iam, req, cred.AccessKey, cred.SecretKey, int64(10*time.Minute.Seconds())); err != nil {
		t.Fatalf("Unable to initialized new signed http request %s", err)
	}
	return req
}

// preSignV4 adds presigned URL parameters to the request
func preSignV4(iam *IdentityAccessManagement, req *http.Request, accessKey, secretKey string, expires int64) error {
	// Create credential scope
	now := time.Now().UTC()
	dateStr := now.Format(iso8601Format)

	// Create credential header
	scope := fmt.Sprintf("%s/%s/%s/%s", now.Format(yyyymmdd), "us-east-1", "s3", "aws4_request")
	credential := fmt.Sprintf("%s/%s", accessKey, scope)

	// Get the query parameters
	query := req.URL.Query()
	query.Set("X-Amz-Algorithm", signV4Algorithm)
	query.Set("X-Amz-Credential", credential)
	query.Set("X-Amz-Date", dateStr)
	query.Set("X-Amz-Expires", fmt.Sprintf("%d", expires))
	query.Set("X-Amz-SignedHeaders", "host")

	// Set the query on the URL (without signature yet)
	req.URL.RawQuery = query.Encode()

	// Get the payload hash
	hashedPayload := getContentSha256Cksum(req)

	// Extract signed headers
	extractedSignedHeaders := make(http.Header)
	extractedSignedHeaders["host"] = []string{req.Host}

	// Get canonical request
	canonicalRequest := getCanonicalRequest(extractedSignedHeaders, hashedPayload, req.URL.RawQuery, req.URL.Path, req.Method)

	// Get string to sign
	stringToSign := getStringToSign(canonicalRequest, now, scope)

	// Get signing key
	signingKey := getSigningKey(secretKey, now.Format(yyyymmdd), "us-east-1", "s3")

	// Calculate signature
	signature := getSignature(signingKey, stringToSign)

	// Add signature to query
	query.Set("X-Amz-Signature", signature)
	req.URL.RawQuery = query.Encode()

	return nil
}

// newTestIAM creates a test IAM with a standard test user
func newTestIAM() *IdentityAccessManagement {
	iam := &IdentityAccessManagement{}
	iam.identities = []*Identity{
		{
			Name:        "testuser",
			Credentials: []*Credential{{AccessKey: "AKIAIOSFODNN7EXAMPLE", SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"}},
			Actions:     []Action{s3_constants.ACTION_ADMIN, s3_constants.ACTION_READ, s3_constants.ACTION_WRITE},
		},
	}
	// Initialize the access key map for lookup
	iam.accessKeyIdent = make(map[string]*Identity)
	iam.accessKeyIdent["AKIAIOSFODNN7EXAMPLE"] = iam.identities[0]
	return iam
}

// Test X-Forwarded-Prefix support for reverse proxy scenarios
func TestSignatureV4WithForwardedPrefix(t *testing.T) {
	tests := []struct {
		name            string
		forwardedPrefix string
		expectedPath    string
	}{
		{
			name:            "prefix without trailing slash",
			forwardedPrefix: "/s3",
			expectedPath:    "/s3/test-bucket/test-object",
		},
		{
			name:            "prefix with trailing slash",
			forwardedPrefix: "/s3/",
			expectedPath:    "/s3/test-bucket/test-object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iam := newTestIAM()

			// Create a request with X-Forwarded-Prefix header
			r, err := newTestRequest("GET", "https://example.com/test-bucket/test-object", 0, nil)
			if err != nil {
				t.Fatalf("Failed to create test request: %v", err)
			}

			// Set the mux variables manually since we're not going through the actual router
			r = mux.SetURLVars(r, map[string]string{
				"bucket": "test-bucket",
				"object": "test-object",
			})

			r.Header.Set("X-Forwarded-Prefix", tt.forwardedPrefix)
			r.Header.Set("Host", "example.com")
			r.Header.Set("X-Forwarded-Host", "example.com")

			// Sign the request with the expected normalized path
			signV4WithPath(r, "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", tt.expectedPath)

			// Test signature verification
			_, errCode := iam.doesSignatureMatch(getContentSha256Cksum(r), r)
			if errCode != s3err.ErrNone {
				t.Errorf("Expected successful signature validation with X-Forwarded-Prefix %q, got error: %v (code: %d)", tt.forwardedPrefix, errCode, int(errCode))
			}
		})
	}
}

// Test X-Forwarded-Port support for reverse proxy scenarios
func TestSignatureV4WithForwardedPort(t *testing.T) {
	tests := []struct {
		name           string
		host           string
		forwardedHost  string
		forwardedPort  string
		forwardedProto string
		expectedHost   string
	}{
		{
			name:           "HTTP with non-standard port",
			host:           "backend:8333",
			forwardedHost:  "example.com",
			forwardedPort:  "8080",
			forwardedProto: "http",
			expectedHost:   "example.com:8080",
		},
		{
			name:           "HTTPS with non-standard port",
			host:           "backend:8333",
			forwardedHost:  "example.com",
			forwardedPort:  "8443",
			forwardedProto: "https",
			expectedHost:   "example.com:8443",
		},
		{
			name:           "HTTP with standard port (80)",
			host:           "backend:8333",
			forwardedHost:  "example.com",
			forwardedPort:  "80",
			forwardedProto: "http",
			expectedHost:   "example.com",
		},
		{
			name:           "HTTPS with standard port (443)",
			host:           "backend:8333",
			forwardedHost:  "example.com",
			forwardedPort:  "443",
			forwardedProto: "https",
			expectedHost:   "example.com",
		},
		{
			name:           "empty proto with non-standard port",
			host:           "backend:8333",
			forwardedHost:  "example.com",
			forwardedPort:  "8080",
			forwardedProto: "",
			expectedHost:   "example.com:8080",
		},
		{
			name:           "empty proto with standard http port",
			host:           "backend:8333",
			forwardedHost:  "example.com",
			forwardedPort:  "80",
			forwardedProto: "",
			expectedHost:   "example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iam := newTestIAM()

			// Create a request
			r, err := newTestRequest("GET", "https://"+tt.host+"/test-bucket/test-object", 0, nil)
			if err != nil {
				t.Fatalf("Failed to create test request: %v", err)
			}

			// Set the mux variables manually since we're not going through the actual router
			r = mux.SetURLVars(r, map[string]string{
				"bucket": "test-bucket",
				"object": "test-object",
			})

			// Set forwarded headers
			r.Header.Set("Host", tt.host)
			r.Header.Set("X-Forwarded-Host", tt.forwardedHost)
			r.Header.Set("X-Forwarded-Port", tt.forwardedPort)
			r.Header.Set("X-Forwarded-Proto", tt.forwardedProto)

			// Sign the request with the expected host header
			// We need to temporarily modify the Host header for signing
            signV4WithPath(r, "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", r.URL.Path)

			// Test signature verification
			_, errCode := iam.doesSignatureMatch(getContentSha256Cksum(r), r)
			if errCode != s3err.ErrNone {
				t.Errorf("Expected successful signature validation with forwarded port, got error: %v (code: %d)", errCode, int(errCode))
			}
		})
	}
}

// Test basic presigned URL functionality without prefix
func TestPresignedSignatureV4Basic(t *testing.T) {
	iam := newTestIAM()

	// Create a presigned request without X-Forwarded-Prefix header
	r, err := newTestRequest("GET", "https://example.com/test-bucket/test-object", 0, nil)
	if err != nil {
		t.Fatalf("Failed to create test request: %v", err)
	}

	// Set the mux variables manually since we're not going through the actual router
	r = mux.SetURLVars(r, map[string]string{
		"bucket": "test-bucket",
		"object": "test-object",
	})

	r.Header.Set("Host", "example.com")

	// Create presigned URL with the normal path (no prefix)
	err = preSignV4WithPath(iam, r, "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", 3600, r.URL.Path)
	if err != nil {
		t.Errorf("Failed to presign request: %v", err)
	}

	// Test presigned signature verification
	_, errCode := iam.doesPresignedSignatureMatch(getContentSha256Cksum(r), r)
	if errCode != s3err.ErrNone {
		t.Errorf("Expected successful presigned signature validation, got error: %v (code: %d)", errCode, int(errCode))
	}
}

// Test X-Forwarded-Prefix support for presigned URLs
func TestPresignedSignatureV4WithForwardedPrefix(t *testing.T) {
	tests := []struct {
		name            string
		forwardedPrefix string
		originalPath    string
		expectedPath    string
	}{
		{
			name:            "prefix without trailing slash",
			forwardedPrefix: "/s3",
			originalPath:    "/s3/test-bucket/test-object",
			expectedPath:    "/s3/test-bucket/test-object",
		},
		{
			name:            "prefix with trailing slash",
			forwardedPrefix: "/s3/",
			originalPath:    "/s3/test-bucket/test-object",
			expectedPath:    "/s3/test-bucket/test-object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			iam := newTestIAM()

			// Create a presigned request that simulates reverse proxy scenario:
			// 1. Client generates presigned URL with prefixed path
			// 2. Proxy strips prefix and forwards to SeaweedFS with X-Forwarded-Prefix header

			// Start with the original request URL (what client sees)
			r, err := newTestRequest("GET", "https://example.com"+tt.originalPath, 0, nil)
			if err != nil {
				t.Fatalf("Failed to create test request: %v", err)
			}

			// Generate presigned URL with the original prefixed path
			err = preSignV4WithPath(iam, r, "AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", 3600, tt.originalPath)
			if err != nil {
				t.Errorf("Failed to presign request: %v", err)
				return
			}

			// Now simulate what the reverse proxy does:
			// 1. Strip the prefix from the URL path
			r.URL.Path = "/test-bucket/test-object"

			// 2. Set the mux variables for the stripped path
			r = mux.SetURLVars(r, map[string]string{
				"bucket": "test-bucket",
				"object": "test-object",
			})

			// 3. Add the forwarded headers
			r.Header.Set("X-Forwarded-Prefix", tt.forwardedPrefix)
			r.Header.Set("Host", "example.com")
			r.Header.Set("X-Forwarded-Host", "example.com")

			// Test presigned signature verification
			_, errCode := iam.doesPresignedSignatureMatch(getContentSha256Cksum(r), r)
			if errCode != s3err.ErrNone {
				t.Errorf("Expected successful presigned signature validation with X-Forwarded-Prefix %q, got error: %v (code: %d)", tt.forwardedPrefix, errCode, int(errCode))
			}
		})
	}
}

// preSignV4WithPath adds presigned URL parameters to the request with a custom path
func preSignV4WithPath(iam *IdentityAccessManagement, req *http.Request, accessKey, secretKey string, expires int64, urlPath string) error {
	// Create credential scope
	now := time.Now().UTC()
	dateStr := now.Format(iso8601Format)

	// Create credential header
	scope := fmt.Sprintf("%s/%s/%s/%s", now.Format(yyyymmdd), "us-east-1", "s3", "aws4_request")
	credential := fmt.Sprintf("%s/%s", accessKey, scope)

	// Get the query parameters
	query := req.URL.Query()
	query.Set("X-Amz-Algorithm", signV4Algorithm)
	query.Set("X-Amz-Credential", credential)
	query.Set("X-Amz-Date", dateStr)
	query.Set("X-Amz-Expires", fmt.Sprintf("%d", expires))
	query.Set("X-Amz-SignedHeaders", "host")

	// Set the query on the URL (without signature yet)
	req.URL.RawQuery = query.Encode()

	// Get the payload hash
	hashedPayload := getContentSha256Cksum(req)

	// Extract signed headers
	extractedSignedHeaders := make(http.Header)
	extractedSignedHeaders["host"] = []string{extractHostHeader(req)}

	// Get canonical request with custom path
	canonicalRequest := getCanonicalRequest(extractedSignedHeaders, hashedPayload, req.URL.RawQuery, urlPath, req.Method)

	// Get string to sign
	stringToSign := getStringToSign(canonicalRequest, now, scope)

	// Get signing key
	signingKey := getSigningKey(secretKey, now.Format(yyyymmdd), "us-east-1", "s3")

	// Calculate signature
	signature := getSignature(signingKey, stringToSign)

	// Add signature to query
	query.Set("X-Amz-Signature", signature)
	req.URL.RawQuery = query.Encode()

	return nil
}

// signV4WithPath signs a request with a custom path
func signV4WithPath(req *http.Request, accessKey, secretKey, urlPath string) {
	// Create credential scope
	now := time.Now().UTC()
	dateStr := now.Format(iso8601Format)

	// Set required headers
	req.Header.Set("X-Amz-Date", dateStr)

	// Create credential header
	scope := fmt.Sprintf("%s/%s/%s/%s", now.Format(yyyymmdd), "us-east-1", "s3", "aws4_request")
	credential := fmt.Sprintf("%s/%s", accessKey, scope)

	// Get signed headers
	signedHeaders := "host;x-amz-date"

	// Extract signed headers
	extractedSignedHeaders := make(http.Header)
	extractedSignedHeaders["host"] = []string{extractHostHeader(req)}
	extractedSignedHeaders["x-amz-date"] = []string{dateStr}

	// Get the payload hash
	hashedPayload := getContentSha256Cksum(req)

	// Get canonical request with custom path
	canonicalRequest := getCanonicalRequest(extractedSignedHeaders, hashedPayload, req.URL.RawQuery, urlPath, req.Method)

	// Get string to sign
	stringToSign := getStringToSign(canonicalRequest, now, scope)

	// Get signing key
	signingKey := getSigningKey(secretKey, now.Format(yyyymmdd), "us-east-1", "s3")

	// Calculate signature
	signature := getSignature(signingKey, stringToSign)

	// Set Authorization header
	authorization := fmt.Sprintf("%s Credential=%s, SignedHeaders=%s, Signature=%s",
		signV4Algorithm, credential, signedHeaders, signature)
	req.Header.Set("Authorization", authorization)
}

// Returns new HTTP request object.
func newTestRequest(method, urlStr string, contentLength int64, body io.ReadSeeker) (*http.Request, error) {
	if method == "" {
		method = http.MethodPost
	}

	// Save for subsequent use
	var hashedPayload string
	var md5Base64 string
	switch {
	case body == nil:
		hashedPayload = getSHA256Hash([]byte{})
	default:
		payloadBytes, err := io.ReadAll(body)
		if err != nil {
			return nil, err
		}
		hashedPayload = getSHA256Hash(payloadBytes)
		md5Base64 = getMD5HashBase64(payloadBytes)
	}
	// Seek back to beginning.
	if body != nil {
		body.Seek(0, 0)
	} else {
		body = bytes.NewReader([]byte(""))
	}
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	if md5Base64 != "" {
		req.Header.Set("Content-Md5", md5Base64)
	}
	req.Header.Set("x-amz-content-sha256", hashedPayload)

	// Add Content-Length
	req.ContentLength = contentLength

	return req, nil
}

// getMD5HashBase64 returns MD5 hash in base64 encoding of given data.
func getMD5HashBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(getMD5Sum(data))
}

// getSHA256Sum returns SHA-256 sum of given data.
func getSHA256Sum(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// getMD5Sum returns MD5 sum of given data.
func getMD5Sum(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// getMD5Hash returns MD5 hash in hex encoding of given data.
func getMD5Hash(data []byte) string {
	return hex.EncodeToString(getMD5Sum(data))
}

var ignoredHeaders = map[string]bool{
	"Authorization":  true,
	"Content-Type":   true,
	"Content-Length": true,
	"User-Agent":     true,
}

// Tests the test helper with an example from the AWS Doc.
// https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html
// This time it's a PUT request uploading the file with content "Welcome to Amazon S3."
func TestGetStringToSignPUT(t *testing.T) {

	canonicalRequest := `PUT
/test%24file.text

date:Fri, 24 May 2013 00:00:00 GMT
host:examplebucket.s3.amazonaws.com
x-amz-content-sha256:44ce7dd67c959e0d3524ffac1771dfbba87d2b6b4b4e99e42034a8b803f8b072
x-amz-date:20130524T000000Z
x-amz-storage-class:REDUCED_REDUNDANCY

date;host;x-amz-content-sha256;x-amz-date;x-amz-storage-class
44ce7dd67c959e0d3524ffac1771dfbba87d2b6b4b4e99e42034a8b803f8b072`

	date, err := time.Parse(iso8601Format, "20130524T000000Z")

	if err != nil {
		t.Fatalf("Error parsing date: %v", err)
	}

	scope := "20130524/us-east-1/s3/aws4_request"
	stringToSign := getStringToSign(canonicalRequest, date, scope)

	expected := `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
9e0e90d9c76de8fa5b200d8c849cd5b8dc7a3be3951ddb7f6a76b4158342019d`

	assert.Equal(t, expected, stringToSign)
}

// Tests the test helper with an example from the AWS Doc.
// https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html
// The GET request example with empty string hash.
func TestGetStringToSignGETEmptyStringHash(t *testing.T) {

	canonicalRequest := `GET
/test.txt

host:examplebucket.s3.amazonaws.com
range:bytes=0-9
x-amz-content-sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
x-amz-date:20130524T000000Z

host;range;x-amz-content-sha256;x-amz-date
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855`

	date, err := time.Parse(iso8601Format, "20130524T000000Z")

	if err != nil {
		t.Fatalf("Error parsing date: %v", err)
	}

	scope := "20130524/us-east-1/s3/aws4_request"
	stringToSign := getStringToSign(canonicalRequest, date, scope)

	expected := `AWS4-HMAC-SHA256
20130524T000000Z
20130524/us-east-1/s3/aws4_request
7344ae5b7ee6c3e7e6b0fe0640412a37625d1fbfff95c48bbb2dc43964946972`

	assert.Equal(t, expected, stringToSign)
}

// Sign given request using Signature V4.
func signRequestV4(req *http.Request, accessKey, secretKey string) error {
	// Get hashed payload.
	hashedPayload := req.Header.Get("x-amz-content-sha256")
	if hashedPayload == "" {
		return fmt.Errorf("Invalid hashed payload")
	}

	currTime := time.Now()

	// Set x-amz-date.
	req.Header.Set("x-amz-date", currTime.Format(iso8601Format))

	// Get header map.
	headerMap := make(map[string][]string)
	for k, vv := range req.Header {
		// If request header key is not in ignored headers, then add it.
		if _, ok := ignoredHeaders[http.CanonicalHeaderKey(k)]; !ok {
			headerMap[strings.ToLower(k)] = vv
		}
	}

	// Get header keys.
	headers := []string{"host"}
	for k := range headerMap {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	region := "us-east-1"

	// Get canonical headers.
	var buf bytes.Buffer
	for _, k := range headers {
		buf.WriteString(k)
		buf.WriteByte(':')
		switch {
		case k == "host":
			buf.WriteString(req.URL.Host)
			fallthrough
		default:
			for idx, v := range headerMap[k] {
				if idx > 0 {
					buf.WriteByte(',')
				}
				buf.WriteString(v)
			}
			buf.WriteByte('\n')
		}
	}
	canonicalHeaders := buf.String()

	// Get signed headers.
	signedHeaders := strings.Join(headers, ";")

	// Get canonical query string.
	req.URL.RawQuery = strings.Replace(req.URL.Query().Encode(), "+", "%20", -1)

	// Get canonical URI.
	canonicalURI := EncodePath(req.URL.Path)

	// Get canonical request.
	// canonicalRequest =
	//  <HTTPMethod>\n
	//  <CanonicalURI>\n
	//  <CanonicalQueryString>\n
	//  <CanonicalHeaders>\n
	//  <SignedHeaders>\n
	//  <HashedPayload>
	//
	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI,
		req.URL.RawQuery,
		canonicalHeaders,
		signedHeaders,
		hashedPayload,
	}, "\n")

	// Get scope.
	scope := strings.Join([]string{
		currTime.Format(yyyymmdd),
		region,
		"s3",
		"aws4_request",
	}, "/")

	stringToSign := "AWS4-HMAC-SHA256" + "\n" + currTime.Format(iso8601Format) + "\n"
	stringToSign = stringToSign + scope + "\n"
	stringToSign = stringToSign + getSHA256Hash([]byte(canonicalRequest))

	date := sumHMAC([]byte("AWS4"+secretKey), []byte(currTime.Format(yyyymmdd)))
	regionHMAC := sumHMAC(date, []byte(region))
	service := sumHMAC(regionHMAC, []byte("s3"))
	signingKey := sumHMAC(service, []byte("aws4_request"))

	signature := hex.EncodeToString(sumHMAC(signingKey, []byte(stringToSign)))

	// final Authorization header
	parts := []string{
		"AWS4-HMAC-SHA256" + " Credential=" + accessKey + "/" + scope,
		"SignedHeaders=" + signedHeaders,
		"Signature=" + signature,
	}
	auth := strings.Join(parts, ", ")
	req.Header.Set("Authorization", auth)

	return nil
}

// EncodePath encode the strings from UTF-8 byte representations to HTML hex escape sequences
//
// This is necessary since regular url.Parse() and url.Encode() functions do not support UTF-8
// non english characters cannot be parsed due to the nature in which url.Encode() is written
//
// This function on the other hand is a direct replacement for url.Encode() technique to support
// pretty much every UTF-8 character.
func EncodePath(pathName string) string {
	if reservedObjectNames.MatchString(pathName) {
		return pathName
	}
	var encodedPathname string
	for _, s := range pathName {
		if 'A' <= s && s <= 'Z' || 'a' <= s && s <= 'z' || '0' <= s && s <= '9' { // §2.3 Unreserved characters (mark)
			encodedPathname = encodedPathname + string(s)
			continue
		}
		switch s {
		case '-', '_', '.', '~', '/': // §2.3 Unreserved characters (mark)
			encodedPathname = encodedPathname + string(s)
			continue
		default:
			runeLen := utf8.RuneLen(s)
			if runeLen < 0 {
				// if utf8 cannot convert return the same string as is
				return pathName
			}
			u := make([]byte, runeLen)
			utf8.EncodeRune(u, s)
			for _, r := range u {
				hex := hex.EncodeToString([]byte{r})
				encodedPathname = encodedPathname + "%" + strings.ToUpper(hex)
			}
		}
	}
	return encodedPathname
}
