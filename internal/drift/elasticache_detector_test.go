package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

func elasticacheStateResource(clusterID, nodeType string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_elasticache_cluster",
		Name: clusterID,
		Attributes: map[string]interface{}{
			"cluster_id": clusterID,
			"node_type":  nodeType,
		},
	}
}

func TestDetectElastiCacheDrift_Missing(t *testing.T) {
	resources := []tfstate.Resource{elasticacheStateResource("my-cluster", "cache.t3.micro")}
	live := []aws.ElastiCacheCluster{}

	drifts := DetectElastiCacheDrift(resources, live)
	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "MISSING")
	assert.Contains(t, drifts[0], "my-cluster")
}

func TestDetectElastiCacheDrift_Deleted(t *testing.T) {
	resources := []tfstate.Resource{elasticacheStateResource("my-cluster", "cache.t3.micro")}
	live := []aws.ElastiCacheCluster{
		{ClusterID: "my-cluster", Status: "deleting", NodeType: "cache.t3.micro"},
	}

	drifts := DetectElastiCacheDrift(resources, live)
	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "DELETED")
}

func TestDetectElastiCacheDrift_NodeTypeMismatch(t *testing.T) {
	resources := []tfstate.Resource{elasticacheStateResource("my-cluster", "cache.t3.micro")}
	live := []aws.ElastiCacheCluster{
		{ClusterID: "my-cluster", Status: "available", NodeType: "cache.t3.medium"},
	}

	drifts := DetectElastiCacheDrift(resources, live)
	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "CHANGED")
	assert.Contains(t, drifts[0], "node_type")
}

func TestDetectElastiCacheDrift_NoDrift(t *testing.T) {
	resources := []tfstate.Resource{elasticacheStateResource("my-cluster", "cache.t3.micro")}
	live := []aws.ElastiCacheCluster{
		{ClusterID: "my-cluster", Status: "available", NodeType: "cache.t3.micro"},
	}

	drifts := DetectElastiCacheDrift(resources, live)
	assert.Empty(t, drifts)
}
