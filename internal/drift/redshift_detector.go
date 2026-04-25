package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
)

type RedshiftDriftResult struct {
	ClusterID string
	Issue     string
}

// DetectRedshiftDrift compares Terraform state resources against live Redshift clusters.
// It reports clusters that are missing, deleted, or have node type mismatches.
func DetectRedshiftDrift(stateResources []StateResource, live []aws.RedshiftCluster) []RedshiftDriftResult {
	liveMap := make(map[string]aws.RedshiftCluster, len(live))
	for _, c := range live {
		liveMap[c.ID] = c
	}

	var results []RedshiftDriftResult

	for _, sr := range stateResources {
		if sr.Type != "aws_redshift_cluster" {
			continue
		}

		clusterID, _ := sr.Attributes["cluster_identifier"].(string)
		if clusterID == "" {
			continue
		}

		liveCluster, found := liveMap[clusterID]
		if !found {
			results = append(results, RedshiftDriftResult{
				ClusterID: clusterID,
				Issue:     "cluster not found in AWS",
			})
			continue
		}

		if liveCluster.Status == "deleted" || liveCluster.Status == "deleting" {
			results = append(results, RedshiftDriftResult{
				ClusterID: clusterID,
				Issue:     fmt.Sprintf("cluster status is '%s'", liveCluster.Status),
			})
			continue
		}

		expectedNodeType, _ := sr.Attributes["node_type"].(string)
		if expectedNodeType != "" && liveCluster.NodeType != expectedNodeType {
			results = append(results, RedshiftDriftResult{
				ClusterID: clusterID,
				Issue:     fmt.Sprintf("node_type mismatch: state=%s live=%s", expectedNodeType, liveCluster.NodeType),
			})
		}
	}

	return results
}
