package drift

import (
	"github.com/your-org/driftctl-lite/internal/aws"
)

type IAMDriftResult struct {
	ResourceID string
	Status     string // "missing" | "ok"
}

// DetectIAMDrift compares Terraform state role IDs against live IAM roles.
// It returns a result for each state resource with a valid ID, marking it
// as "missing" if no matching live IAM role is found, or "ok" otherwise.
func DetectIAMDrift(stateResources []map[string]interface{}, liveRoles []aws.IAMRole) []IAMDriftResult {
	liveIndex := make(map[string]bool, len(liveRoles))
	for _, r := range liveRoles {
		liveIndex[r.ID] = true
	}

	var results []IAMDriftResult
	for _, res := range stateResources {
		id, ok := res["id"].(string)
		if !ok || id == "" {
			continue
		}
		status := "ok"
		if !liveIndex[id] {
			status = "missing"
		}
		results = append(results, IAMDriftResult{ResourceID: id, Status: status})
	}
	return results
}

// FilterByStatus returns only the IAMDriftResults that match the given status.
// Valid status values are "ok" and "missing".
func FilterByStatus(results []IAMDriftResult, status string) []IAMDriftResult {
	var filtered []IAMDriftResult
	for _, r := range results {
		if r.Status == status {
			filtered = append(filtered, r)
		}
	}
	return filtered
}
