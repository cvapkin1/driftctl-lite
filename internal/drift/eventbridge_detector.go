package drift

import (
	"fmt"

	"github.com/owner/driftctl-lite/internal/aws"
	"github.com/owner/driftctl-lite/internal/tfstate"
)

// DetectEventBridgeDrift compares EventBridge rules from Terraform state
// against live AWS resources and returns a list of drift messages.
func DetectEventBridgeDrift(resources []tfstate.Resource, live []aws.EventBridgeRule) []string {
	var drifts []string

	liveIndex := make(map[string]aws.EventBridgeRule, len(live))
	for _, r := range live {
		liveIndex[r.Name] = r
	}

	for _, res := range resources {
		if res.Type != "aws_cloudwatch_event_rule" {
			continue
		}

		name, _ := res.Attributes["name"].(string)
		if name == "" {
			continue
		}

		liveRule, found := liveIndex[name]
		if !found {
			drifts = append(drifts, fmt.Sprintf("[MISSING] EventBridge rule %q not found in AWS", name))
			continue
		}

		if liveRule.State == "DISABLED" {
			drifts = append(drifts, fmt.Sprintf("[DISABLED] EventBridge rule %q is disabled in AWS", name))
		}

		expectedARN, _ := res.Attributes["arn"].(string)
		if expectedARN != "" && liveRule.ARN != expectedARN {
			drifts = append(drifts, fmt.Sprintf(
				"[ARN_MISMATCH] EventBridge rule %q: state has %q, live has %q",
				name, expectedARN, liveRule.ARN,
			))
		}
	}

	return drifts
}
