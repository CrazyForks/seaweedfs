package topic

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/seaweedfs/seaweedfs/weed/pb/mq_pb"
	"github.com/shirou/gopsutil/v3/cpu"
	"time"
)

// LocalTopicManager manages topics on local broker
type LocalTopicManager struct {
	topics cmap.ConcurrentMap[string, *LocalTopic]
}

// NewLocalTopicManager creates a new LocalTopicManager
func NewLocalTopicManager() *LocalTopicManager {
	return &LocalTopicManager{
		topics: cmap.New[*LocalTopic](),
	}
}

// AddTopic adds a topic to the local topic manager
func (manager *LocalTopicManager) AddTopicPartition(topic Topic, localPartition *LocalPartition) {
	localTopic, ok := manager.topics.Get(topic.String())
	if !ok {
		localTopic = &LocalTopic{
			Topic:      topic,
			Partitions: make([]*LocalPartition, 0),
		}
	}
	manager.topics.SetIfAbsent(topic.String(), localTopic)
	if localTopic.findPartition(localPartition.Partition) != nil {
		return
	}
	localTopic.Partitions = append(localTopic.Partitions, localPartition)
}

// GetTopic gets a topic from the local topic manager
func (manager *LocalTopicManager) GetTopicPartition(topic Topic, partition Partition) *LocalPartition {
	localTopic, ok := manager.topics.Get(topic.String())
	if !ok {
		return nil
	}
	return localTopic.findPartition(partition)
}

// RemoveTopic removes a topic from the local topic manager
func (manager *LocalTopicManager) RemoveTopic(topic Topic) {
	manager.topics.Remove(topic.String())
}

func (manager *LocalTopicManager) RemoveTopicPartition(topic Topic, partition Partition) (removed bool) {
	localTopic, ok := manager.topics.Get(topic.String())
	if !ok {
		return false
	}
	return localTopic.removePartition(partition)
}

func (manager *LocalTopicManager) CollectStats(duration time.Duration) *mq_pb.BrokerStats {
	stats := &mq_pb.BrokerStats{
		Stats: make(map[string]*mq_pb.TopicPartitionStats),
	}
	manager.topics.IterCb(func(topic string, localTopic *LocalTopic) {
		for _, localPartition := range localTopic.Partitions {
			stats.Stats[topic] = &mq_pb.TopicPartitionStats{
				Topic: &mq_pb.Topic{
					Namespace: string(localTopic.Namespace),
					Name:      localTopic.Name,
				},
				Partition: &mq_pb.Partition{
					RingSize:   localPartition.RingSize,
					RangeStart: localPartition.RangeStart,
					RangeStop:  localPartition.RangeStop,
				},
				ConsumerCount: localPartition.ConsumerCount,
			}
		}
	})

	// collect current broker's cpu usage
	usages, err := cpu.Percent(duration, false)
	if err == nil && len(usages) > 0 {
		stats.CpuUsagePercent = int32(usages[0])
	}

	return stats

}
