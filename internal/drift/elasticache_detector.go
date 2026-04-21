package drift

import (
	"fmt"

	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

// DetectElastiCacheDrift compares ElastiCache clusters in Terraform state
// against live AWS resources and returns a list of drift messages.
func DetectElastiCacheDrift(resources []tfstate.Resource, live []aws.ElastiCacheCluster) []string {
	liveMap := make(map[string]aws.ElastiCacheCluster)
	for _, c := range live {
		liveMap[c.ClusterID] = c
	}

	var drifts []string
	for _, r := range resources {
		if r.Type != "aws_elasticache_cluster" {
			continue
		}

		clusterID, _ := r.Attributes["cluster_id"].(string)
		liveCluster, found := liveMap[clusterID]
		if !found {
			drifts = append(drifts, fmt.Sprintf("[MISSING] ElastiCache cluster %q not found in AWS", clusterID))
			continue
		}

		if liveCluster.Status == "deleting" || liveCluster.Status == "deleted" {
			drifts = append(drifts, fmt.Sprintf("[DELETED] ElastiCache cluster %q has status %q", clusterID, liveCluster.Status))
			continue
		}

		if stateNodeType, ok := r.Attributes["node_type"].(string); ok {
			if stateNodeType != liveCluster.NodeType {
				drifts = append(drifts, fmt.Sprintf(
					"[CHANGED] ElastiCache cluster %q node_type: state=%q live=%q",
					clusterID, stateNodeType, liveCluster.NodeType,
				))
			}
		}
	}

	return drifts
}
