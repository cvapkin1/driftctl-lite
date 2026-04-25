package drift

import (
	"fmt"

	"github.com/elie-h/driftctl-lite/internal/aws"
)

type EMRDriftResult struct {
	ClusterID string
	Status    string
	Message   string
}

// DetectEMRDrift compares EMR clusters in Terraform state against live AWS resources.
// It reports missing clusters and clusters in a terminated or terminating state.
func DetectEMRDrift(stateResources []map[string]interface{}, live []aws.EMRCluster) []EMRDriftResult {
	var results []EMRDriftResult

	liveByID := make(map[string]aws.EMRCluster, len(live))
	for _, c := range live {
		liveByID[c.ID] = c
	}

	for _, res := range stateResources {
		id, _ := res["id"].(string)
		if id == "" {
			continue
		}

		liveCluster, found := liveByID[id]
		if !found {
			results = append(results, EMRDriftResult{
				ClusterID: id,
				Status:    "missing",
				Message:   fmt.Sprintf("EMR cluster %s not found in AWS", id),
			})
			continue
		}

		switch liveCluster.State {
		case "TERMINATED", "TERMINATING", "TERMINATED_WITH_ERRORS":
			results = append(results, EMRDriftResult{
				ClusterID: id,
				Status:    "terminated",
				Message:   fmt.Sprintf("EMR cluster %s is in state %s", id, liveCluster.State),
			})
		}
	}

	return results
}
