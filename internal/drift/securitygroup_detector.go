package drift

import (
	"fmt"

	"github.com/owner/driftctl-lite/internal/aws"
	"github.com/owner/driftctl-lite/internal/tfstate"
)

// DetectSecurityGroupDrift compares Terraform state security group resources
// against live AWS security groups and reports any drift.
func DetectSecurityGroupDrift(resources []tfstate.Resource, live []aws.SecurityGroup) []DriftResult {
	liveMap := make(map[string]aws.SecurityGroup)
	for _, sg := range live {
		liveMap[sg.ID] = sg
	}

	var results []DriftResult

	for _, res := range resources {
		if res.Type != "aws_security_group" {
			continue
		}

		sgID, _ := res.Attributes["id"].(string)
		if sgID == "" {
			continue
		}

		live, found := liveMap[sgID]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   sgID,
				ResourceType: res.Type,
				DriftType:    "missing",
				Details:      fmt.Sprintf("security group %s not found in AWS", sgID),
			})
			continue
		}

		stateDesc, _ := res.Attributes["description"].(string)
		if stateDesc != "" && live.Description != stateDesc {
			results = append(results, DriftResult{
				ResourceID:   sgID,
				ResourceType: res.Type,
				DriftType:    "description_mismatch",
				Details:      fmt.Sprintf("expected description %q, got %q", stateDesc, live.Description),
			})
		}

		stateVPC, _ := res.Attributes["vpc_id"].(string)
		if stateVPC != "" && live.VPCID != stateVPC {
			results = append(results, DriftResult{
				ResourceID:   sgID,
				ResourceType: res.Type,
				DriftType:    "vpc_mismatch",
				Details:      fmt.Sprintf("expected vpc_id %q, got %q", stateVPC, live.VPCID),
			})
		}
	}

	return results
}
