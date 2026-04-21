package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	elasticachetypes "github.com/aws/aws-sdk-go-v2/service/elasticache/types"
	"github.com/stretchr/testify/assert"
)

func TestMapElastiCacheCluster_Fields(t *testing.T) {
	input := elasticachetypes.CacheCluster{
		CacheClusterId:     aws.String("my-cluster"),
		Engine:             aws.String("redis"),
		CacheClusterStatus: aws.String("available"),
		CacheNodeType:      aws.String("cache.t3.micro"),
	}

	result := mapElastiCacheCluster(input)

	assert.Equal(t, "my-cluster", result.ClusterID)
	assert.Equal(t, "redis", result.Engine)
	assert.Equal(t, "available", result.Status)
	assert.Equal(t, "cache.t3.micro", result.NodeType)
}

func TestMapElastiCacheCluster_EmptyFields(t *testing.T) {
	input := elasticachetypes.CacheCluster{}
	result := mapElastiCacheCluster(input)
	assert.Equal(t, "", result.ClusterID)
	assert.Equal(t, "", result.Engine)
	assert.Equal(t, "", result.Status)
	assert.Equal(t, "", result.NodeType)
}
