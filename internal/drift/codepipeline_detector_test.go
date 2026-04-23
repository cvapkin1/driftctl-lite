package drift

import (
	"testing"

	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

func codepipelineStateResource(name, arn string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_codepipeline",
		Name: name,
		Attributes: map[string]interface{}{
			"name": name,
			"arn":  arn,
		},
	}
}

func TestDetectCodePipelineDrift_Missing(t *testing.T) {
	state := []tfstate.Resource{codepipelineStateResource("my-pipeline", "arn:aws:codepipeline:us-east-1:123:my-pipeline")}
	live := []aws.CodePipelineResource{}

	results := DetectCodePipelineDrift(state, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].DriftType != "MISSING" {
		t.Errorf("expected MISSING, got %s", results[0].DriftType)
	}
}

func TestDetectCodePipelineDrift_ARNMismatch(t *testing.T) {
	state := []tfstate.Resource{codepipelineStateResource("my-pipeline", "arn:aws:codepipeline:us-east-1:123:my-pipeline")}
	live := []aws.CodePipelineResource{
		{Name: "my-pipeline", ARN: "arn:aws:codepipeline:us-east-1:999:my-pipeline", Status: "ACTIVE"},
	}

	results := DetectCodePipelineDrift(state, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].DriftType != "ARN_MISMATCH" {
		t.Errorf("expected ARN_MISMATCH, got %s", results[0].DriftType)
	}
}

func TestDetectCodePipelineDrift_NoDrift(t *testing.T) {
	arn := "arn:aws:codepipeline:us-east-1:123:my-pipeline"
	state := []tfstate.Resource{codepipelineStateResource("my-pipeline", arn)}
	live := []aws.CodePipelineResource{
		{Name: "my-pipeline", ARN: arn, Status: "ACTIVE"},
	}

	results := DetectCodePipelineDrift(state, live)
	if len(results) != 0 {
		t.Errorf("expected no drift, got %d results", len(results))
	}
}
