package drift

import (
	"fmt"

	"github.com/edobry/driftctl-lite/internal/aws"
	"github.com/edobry/driftctl-lite/internal/tfstate"
)

// DetectSSMDrift compares SSM parameters in Terraform state against live AWS resources.
func DetectSSMDrift(resources []tfstate.Resource, live []aws.SSMParameter) []string {
	var drifts []string

	liveByName := make(map[string]aws.SSMParameter)
	for _, p := range live {
		liveByName[p.Name] = p
	}

	for _, res := range resources {
		if res.Type != "aws_ssm_parameter" {
			continue
		}

		name, _ := res.Attributes["name"].(string)
		if name == "" {
			continue
		}

		liveParam, found := liveByName[name]
		if !found {
			drifts = append(drifts, fmt.Sprintf("MISSING: SSM parameter %q not found in AWS", name))
			continue
		}

		if expectedType, ok := res.Attributes["type"].(string); ok && expectedType != "" {
			if liveParam.Type != expectedType {
				drifts = append(drifts, fmt.Sprintf(
					"TYPE_MISMATCH: SSM parameter %q expected type %q but got %q",
					name, expectedType, liveParam.Type,
				))
			}
		}

		if expectedARN, ok := res.Attributes["arn"].(string); ok && expectedARN != "" {
			if liveParam.ARN != expectedARN {
				drifts = append(drifts, fmt.Sprintf(
					"ARN_MISMATCH: SSM parameter %q expected ARN %q but got %q",
					name, expectedARN, liveParam.ARN,
				))
			}
		}
	}

	return drifts
}
