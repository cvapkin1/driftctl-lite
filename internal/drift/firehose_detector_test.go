package drift

import (
	"testing"

	"github.com/eliraz-refael/driftctl-lite/internal/aws"
	"github.com/eliraz-refael/driftctl-lite/internal/tfstate"
	"github.com/stretchr/testify/assert"
)

func firehoseStateResource(name, arn string) tfstate.Resource {
	return tfstate.Resource{
		Type: firehoseResourceType,
		Name: name,
		Attributes: map[string]interface{}{
			"name": name,
			"arn":  arn,
		},
	}
}

func TestDetectFirehoseDrift_Missing(t *testing.T) {
	res := []tfstate.Resource{firehoseStateResource("my-stream", "arn:aws:firehose:::my-stream")}
	live := []aws.FirehoseDeliveryStream{}

	results := DetectFirehoseDrift(res, live)

	assert.Len(t, results, 1)
	assert.Equal(t, DriftTypeMissing, results[0].DriftType)
	assert.Equal(t, "my-stream", results[0].ResourceID)
}

func TestDetectFirehoseDrift_InactiveStatus(t *testing.T) {
	res := []tfstate.Resource{firehoseStateResource("my-stream", "arn:aws:firehose:::my-stream")}
	live := []aws.FirehoseDeliveryStream{
		{Name: "my-stream", ARN: "arn:aws:firehose:::my-stream", Status: "DELETING"},
	}

	results := DetectFirehoseDrift(res, live)

	assert.Len(t, results, 1)
	assert.Equal(t, DriftTypeModified, results[0].DriftType)
	assert.Contains(t, results[0].Details, "DELETING")
}

func TestDetectFirehoseDrift_ARNMismatch(t *testing.T) {
	res := []tfstate.Resource{firehoseStateResource("my-stream", "arn:aws:firehose:::old-arn")}
	live := []aws.FirehoseDeliveryStream{
		{Name: "my-stream", ARN: "arn:aws:firehose:::new-arn", Status: "ACTIVE"},
	}

	results := DetectFirehoseDrift(res, live)

	assert.Len(t, results, 1)
	assert.Equal(t, DriftTypeModified, results[0].DriftType)
	assert.Contains(t, results[0].Details, "ARN mismatch")
}

func TestDetectFirehoseDrift_NoDrift(t *testing.T) {
	arn := "arn:aws:firehose:us-east-1:123:deliverystream/my-stream"
	res := []tfstate.Resource{firehoseStateResource("my-stream", arn)}
	live := []aws.FirehoseDeliveryStream{
		{Name: "my-stream", ARN: arn, Status: "ACTIVE"},
	}

	results := DetectFirehoseDrift(res, live)

	assert.Empty(t, results)
}
