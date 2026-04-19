package drift

import (
	"github.com/example/driftctl-lite/internal/aws"
	"github.com/example/driftctl-lite/internal/tfstate"
)

// DriftResult holds a single drifted resource.
type DriftResult struct {
	ResourceType string
	ResourceID   string
	Reason       string
}

// DetectEC2Drift compares Terraform state resources against live EC2 instances.
func DetectEC2Drift(resources []tfstate.Resource, live []aws.EC2Instance) []DriftResult {
	liveMap := make(map[string]aws.EC2Instance, len(live))
	for _, inst := range live {
		liveMap[inst.ID] = inst
	}

	var results []DriftResult
	for _, r := range resources {
		if r.Type != "aws_instance" {
			continue
		}
		id, ok := r.Attributes["id"].(string)
		if !ok || id == "" {
			continue
		}
		inst, found := liveMap[id]
		if !found {
			results = append(results, DriftResult{
				ResourceType: r.Type,
				ResourceID:   id,
				Reason:       "resource not found in AWS",
			})
			continue
		}
		if inst.State == "terminated" {
			results = append(results, DriftResult{
				ResourceType: r.Type,
				ResourceID:   id,
				Reason:       "instance is terminated",
			})
		}
	}
	return results
}
