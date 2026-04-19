package drift

import (
	"testing"

	"github.com/edobry/driftctl-lite/internal/aws"
	"github.com/edobry/driftctl-lite/internal/tfstate"
	"github.com/stretchr/testify/assert"
)

func sqsStateResource(name, arn string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_sqs_queue",
		Name: name,
		Attributes: map[string]interface{}{
			"name": name,
			"arn":  arn,
		},
	}
}

func TestDetectSQSDrift_Missing(t *testing.T) {
	res := []tfstate.Resource{sqsStateResource("my-queue", "arn:aws:sqs:us-east-1:123:my-queue")}
	live := []aws.SQSQueue{}
	results := DetectSQSDrift(res, live)
	assert.Len(t, results, 1)
	assert.Equal(t, "missing", results[0].Status)
	assert.Equal(t, "my-queue", results[0].ResourceID)
}

func TestDetectSQSDrift_ARNMismatch(t *testing.T) {
	res := []tfstate.Resource{sqsStateResource("my-queue", "arn:aws:sqs:us-east-1:123:my-queue")}
	live := []aws.SQSQueue{{Name: "my-queue", ARN: "arn:aws:sqs:us-east-1:999:my-queue"}}
	results := DetectSQSDrift(res, live)
	assert.Len(t, results, 1)
	assert.Equal(t, "drifted", results[0].Status)
}

func TestDetectSQSDrift_NoDrift(t *testing.T) {
	arn := "arn:aws:sqs:us-east-1:123:my-queue"
	res := []tfstate.Resource{sqsStateResource("my-queue", arn)}
	live := []aws.SQSQueue{{Name: "my-queue", ARN: arn}}
	results := DetectSQSDrift(res, live)
	assert.Empty(t, results)
}
