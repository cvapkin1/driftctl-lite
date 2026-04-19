package drift

import (
	"fmt"

	"github.com/edobry/driftctl-lite/internal/aws"
	"github.com/edobry/driftctl-lite/internal/tfstate"
)

func DetectSQSDrift(resources []tfstate.Resource, live []aws.SQSQueue) []DriftResult {
	liveMap := make(map[string]aws.SQSQueue)
	for _, q := range live {
		liveMap[q.Name] = q
	}

	var results []DriftResult
	for _, res := range resources {
		if res.Type != "aws_sqs_queue" {
			continue
		}
		name, _ := res.Attributes["name"].(string)
		if name == "" {
			continue
		}
		q, found := liveMap[name]
		if !found {
			results = append(results, DriftResult{
				ResourceType: "aws_sqs_queue",
				ResourceID:   name,
				Status:       "missing",
				Details:      fmt.Sprintf("SQS queue %q not found in AWS", name),
			})
			continue
		}
		expectedARN, _ := res.Attributes["arn"].(string)
		if expectedARN != "" && q.ARN != expectedARN {
			results = append(results, DriftResult{
				ResourceType: "aws_sqs_queue",
				ResourceID:   name,
				Status:       "drifted",
				Details:      fmt.Sprintf("ARN mismatch: state=%s live=%s", expectedARN, q.ARN),
			})
		}
	}
	return results
}
