package drift

import (
	"fmt"

	"github.com/esenac/driftctl-lite/internal/aws"
)

type KinesisDriftResult struct {
	StreamName string
	Status     string
	Details    string
}

// DetectKinesisDrift compares Terraform state resources against live Kinesis streams.
// It reports streams that are missing or in a non-ACTIVE state.
func DetectKinesisDrift(stateResources []map[string]interface{}, live []aws.KinesisStream) []KinesisDriftResult {
	var results []KinesisDriftResult

	liveByARN := make(map[string]aws.KinesisStream, len(live))
	for _, s := range live {
		liveByARN[s.ARN] = s
	}

	for _, res := range stateResources {
		arn, _ := res["arn"].(string)
		name, _ := res["name"].(string)

		liveStream, found := liveByARN[arn]
		if !found {
			results = append(results, KinesisDriftResult{
				StreamName: name,
				Status:     "missing",
				Details:    fmt.Sprintf("stream %q not found in AWS (arn: %s)", name, arn),
			})
			continue
		}

		if liveStream.Status != "ACTIVE" {
			results = append(results, KinesisDriftResult{
				StreamName: name,
				Status:     "inactive",
				Details:    fmt.Sprintf("stream %q has status %q (expected ACTIVE)", name, liveStream.Status),
			})
		}
	}

	return results
}
