package drift

import (
	"fmt"

	"github.com/owner/driftctl-lite/internal/aws"
	"github.com/owner/driftctl-lite/internal/tfstate"
)

// DetectSageMakerDrift compares Terraform state SageMaker models against live AWS models.
func DetectSageMakerDrift(stateResources []tfstate.Resource, liveModels []aws.SageMakerModel) []DriftResult {
	var results []DriftResult

	liveIndex := make(map[string]aws.SageMakerModel, len(liveModels))
	for _, m := range liveModels {
		liveIndex[m.Name] = m
	}

	for _, res := range stateResources {
		if res.Type != "aws_sagemaker_model" {
			continue
		}

		name, _ := res.Attributes["name"].(string)
		expectedARN, _ := res.Attributes["arn"].(string)

		live, found := liveIndex[name]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   name,
				ResourceType: res.Type,
				DriftType:    "missing",
				Details:      fmt.Sprintf("SageMaker model %q not found in AWS", name),
			})
			continue
		}

		if expectedARN != "" && live.ARN != expectedARN {
			results = append(results, DriftResult{
				ResourceID:   name,
				ResourceType: res.Type,
				DriftType:    "arn_mismatch",
				Details:      fmt.Sprintf("ARN mismatch: state=%q live=%q", expectedARN, live.ARN),
			})
		}
	}

	return results
}
