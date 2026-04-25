package drift

import (
	"fmt"

	"github.com/edobry/driftctl-lite/internal/aws"
	"github.com/edobry/driftctl-lite/internal/tfstate"
)

// DetectLightsailDrift compares Terraform state Lightsail instances against live AWS instances.
func DetectLightsailDrift(resources []tfstate.Resource, live []aws.LightsailInstance) []DriftResult {
	var results []DriftResult

	liveMap := make(map[string]aws.LightsailInstance, len(live))
	for _, inst := range live {
		liveMap[inst.Name] = inst
	}

	for _, res := range resources {
		if res.Type != "aws_lightsail_instance" {
			continue
		}

		name, _ := res.Attributes["name"].(string)
		expectedARN, _ := res.Attributes["arn"].(string)

		inst, found := liveMap[name]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   name,
				ResourceType: res.Type,
				DriftType:    DriftTypeMissing,
				Details:      fmt.Sprintf("Lightsail instance %q not found in AWS", name),
			})
			continue
		}

		if inst.State == "stopped" || inst.State == "terminated" {
			results = append(results, DriftResult{
				ResourceID:   name,
				ResourceType: res.Type,
				DriftType:    DriftTypeModified,
				Details:      fmt.Sprintf("Lightsail instance %q is in state %q", name, inst.State),
			})
			continue
		}

		if expectedARN != "" && inst.ARN != expectedARN {
			results = append(results, DriftResult{
				ResourceID:   name,
				ResourceType: res.Type,
				DriftType:    DriftTypeModified,
				Details:      fmt.Sprintf("Lightsail instance %q ARN mismatch: state=%q live=%q", name, expectedARN, inst.ARN),
			})
		}
	}

	return results
}
