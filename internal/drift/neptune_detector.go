package drift

import (
	"fmt"

	"github.com/owner/driftctl-lite/internal/aws"
	"github.com/owner/driftctl-lite/internal/tfstate"
)

// DetectNeptuneDrift compares Neptune clusters in Terraform state against live AWS resources.
// It reports clusters that are missing, deleted, or have a version mismatch.
func DetectNeptuneDrift(resources []tfstate.Resource, live []aws.NeptuneCluster) []DriftResult {
	var results []DriftResult

	liveMap := make(map[string]aws.NeptuneCluster, len(live))
	for _, c := range live {
		liveMap[c.ClusterID] = c
	}

	for _, res := range resources {
		if res.Type != "aws_neptune_cluster" {
			continue
		}

		clusterID, _ := res.Attributes["cluster_identifier"].(string)
		if clusterID == "" {
			continue
		}

		liveCluster, found := liveMap[clusterID]
		if !found {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   clusterID,
				DriftType:    "missing",
				Message:      fmt.Sprintf("Neptune cluster %q not found in AWS", clusterID),
			})
			continue
		}

		if liveCluster.Status == "deleting" || liveCluster.Status == "deleted" {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   clusterID,
				DriftType:    "deleted",
				Message:      fmt.Sprintf("Neptune cluster %q has status %q", clusterID, liveCluster.Status),
			})
			continue
		}

		stateVersion, _ := res.Attributes["engine_version"].(string)
		if stateVersion != "" && stateVersion != liveCluster.EngineVersion {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   clusterID,
				DriftType:    "modified",
				Message:      fmt.Sprintf("Neptune cluster %q engine version mismatch: state=%q live=%q", clusterID, stateVersion, liveCluster.EngineVersion),
			})
		}
	}

	return results
}
