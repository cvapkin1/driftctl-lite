package drift

import (
	"fmt"

	"github.com/owner/driftctl-lite/internal/aws"
)

// RDSDriftResult holds drift information for an RDS instance.
type RDSDriftResult struct {
	InstanceID string
	Status     string
	Message    string
}

// DetectRDSDrift compares Terraform state RDS resources against live AWS RDS instances.
func DetectRDSDrift(stateResources []map[string]interface{}, fetcher *aws.RDSFetcher) ([]RDSDriftResult, error) {
	liveInstances, err := fetcher.FetchAll()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RDS instances: %w", err)
	}

	liveMap := make(map[string]aws.RDSInstance)
	for _, inst := range liveInstances {
		liveMap[inst.DBInstanceIdentifier] = inst
	}

	var results []RDSDriftResult

	for _, res := range stateResources {
		id, _ := res["id"].(string)
		if id == "" {
			continue
		}

		live, found := liveMap[id]
		if !found {
			results = append(results, RDSDriftResult{
				InstanceID: id,
				Status:     "missing",
				Message:    fmt.Sprintf("RDS instance %s not found in AWS", id),
			})
			continue
		}

		if live.DBInstanceStatus == "deleted" || live.DBInstanceStatus == "deleting" {
			results = append(results, RDSDriftResult{
				InstanceID: id,
				Status:     "deleted",
				Message:    fmt.Sprintf("RDS instance %s is in state: %s", id, live.DBInstanceStatus),
			})
			continue
		}

		results = append(results, RDSDriftResult{
			InstanceID: id,
			Status:     "ok",
			Message:    fmt.Sprintf("RDS instance %s is %s", id, live.DBInstanceStatus),
		})
	}

	return results, nil
}
