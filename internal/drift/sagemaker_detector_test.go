package drift

import (
	"testing"

	"github.com/owner/driftctl-lite/internal/aws"
	"github.com/owner/driftctl-lite/internal/tfstate"
	"github.com/stretchr/testify/assert"
)

func sageMakerStateResource(name, arn string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_sagemaker_model",
		Name: name,
		Attributes: map[string]interface{}{
			"name": name,
			"arn":  arn,
		},
	}
}

func TestDetectSageMakerDrift_Missing(t *testing.T) {
	state := []tfstate.Resource{sageMakerStateResource("my-model", "arn:aws:sagemaker:::model/my-model")}
	live := []aws.SageMakerModel{}
	results := DetectSageMakerDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, "missing", results[0].DriftType)
	assert.Equal(t, "my-model", results[0].ResourceID)
}

func TestDetectSageMakerDrift_ARNMismatch(t *testing.T) {
	state := []tfstate.Resource{sageMakerStateResource("my-model", "arn:aws:sagemaker:::model/my-model")}
	live := []aws.SageMakerModel{
		{Name: "my-model", ARN: "arn:aws:sagemaker:::model/different-arn"},
	}
	results := DetectSageMakerDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, "arn_mismatch", results[0].DriftType)
}

func TestDetectSageMakerDrift_NoDrift(t *testing.T) {
	arn := "arn:aws:sagemaker:us-east-1:123:model/my-model"
	state := []tfstate.Resource{sageMakerStateResource("my-model", arn)}
	live := []aws.SageMakerModel{
		{Name: "my-model", ARN: arn},
	}
	results := DetectSageMakerDrift(state, live)
	assert.Empty(t, results)
}
