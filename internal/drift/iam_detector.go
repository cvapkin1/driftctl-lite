package drift

import (
	"github.com/your-org/driftctl-lite/internal/aws"
)

type IAMDriftResult struct {
	ResourceID string
	Status     string // "missing" | "ok"
}

// DetectIAMDrift compares Terraform state role IDs against live IAM roles.
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
