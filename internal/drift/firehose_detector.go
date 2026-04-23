package drift

import (
	"fmt"

	"github.com/eliraz-refael/driftctl-lite/internal/aws"
	"github.com/eliraz-refael/driftctl-lite/internal/tfstate"
)

const (
	firehoseResourceType = "aws_kinesis_firehose_delivery_stream"
	firehoseActiveStatus = "ACTIVE"
)

// DetectFirehoseDrift compares Firehose delivery streams in Terraform state
// against live AWS resources and returns a list of drift results.
func DetectFirehoseDrift(resources []tfstate.Resource, live []aws.FirehoseDeliveryStream) []DriftResult {
	liveMap := make(map[string]aws.FirehoseDeliveryStream)
	for _, s := range live {
		liveMap[s.Name] = s
	}

	var results []DriftResult

	for _, res := range resources {
		if res.Type != firehoseResourceType {
			continue
		}

		name, _ := res.Attributes["name"].(string)
		expectedARN, _ := res.Attributes["arn"].(string)

		liveStream, found := liveMap[name]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   name,
				ResourceType: firehoseResourceType,
				DriftType:    DriftTypeMissing,
				Details:      fmt.Sprintf("delivery stream %q not found in AWS", name),
			})
			continue
		}

		if liveStream.Status != firehoseActiveStatus {
			results = append(results, DriftResult{
				ResourceID:   name,
				ResourceType: firehoseResourceType,
				DriftType:    DriftTypeModified,
				Details:      fmt.Sprintf("delivery stream %q has status %q, expected ACTIVE", name, liveStream.Status),
			})
			continue
		}

		if expectedARN != "" && liveStream.ARN != expectedARN {
			results = append(results, DriftResult{
				ResourceID:   name,
				ResourceType: firehoseResourceType,
				DriftType:    DriftTypeModified,
				Details:      fmt.Sprintf("delivery stream %q ARN mismatch: state=%q live=%q", name, expectedARN, liveStream.ARN),
			})
		}
	}

	return results
}
