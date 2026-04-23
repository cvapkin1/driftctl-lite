package drift

import (
	"fmt"

	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

// DetectCodePipelineDrift compares Terraform state resources against live CodePipeline pipelines.
// It reports pipelines that are missing or have a mismatched ARN.
func DetectCodePipelineDrift(stateResources []tfstate.Resource, live []aws.CodePipelineResource) []DriftResult {
	var results []DriftResult

	liveMap := make(map[string]aws.CodePipelineResource, len(live))
	for _, p := range live {
		liveMap[p.Name] = p
	}

	for _, res := range stateResources {
		if res.Type != "aws_codepipeline" {
			continue
		}

		name, _ := res.Attributes["name"].(string)
		expectedARN, _ := res.Attributes["arn"].(string)

		liveP, found := liveMap[name]
		if !found {
			results = append(results, DriftResult{
				ResourceType: "aws_codepipeline",
				ResourceID:   name,
				DriftType:    "MISSING",
				Detail:       fmt.Sprintf("pipeline %q not found in AWS", name),
			})
			continue
		}

		if expectedARN != "" && liveP.ARN != expectedARN {
			results = append(results, DriftResult{
				ResourceType: "aws_codepipeline",
				ResourceID:   name,
				DriftType:    "ARN_MISMATCH",
				Detail:       fmt.Sprintf("expected ARN %q, got %q", expectedARN, liveP.ARN),
			})
		}
	}

	return results
}
