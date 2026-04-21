package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

// DetectStepFunctionsDrift compares Terraform state Step Functions state machines
// against live AWS state machines and reports any drift.
func DetectStepFunctionsDrift(resources []tfstate.Resource, live []aws.StateMachine) []string {
	var drifts []string

	liveMap := make(map[string]aws.StateMachine, len(live))
	for _, sm := range live {
		liveMap[sm.ARN] = sm
	}

	for _, res := range resources {
		if res.Type != "aws_sfn_state_machine" {
			continue
		}

		arn, _ := res.Attributes["arn"].(string)
		name, _ := res.Attributes["name"].(string)

		liveSM, found := liveMap[arn]
		if !found {
			drifts = append(drifts, fmt.Sprintf("[MISSING] aws_sfn_state_machine %q (arn: %s) not found in AWS", name, arn))
			continue
		}

		if liveSM.Status == "DELETING" || liveSM.Status == "DELETED" {
			drifts = append(drifts, fmt.Sprintf("[DELETED] aws_sfn_state_machine %q is in status %q in AWS", name, liveSM.Status))
			continue
		}

		if liveSM.Name != name {
			drifts = append(drifts, fmt.Sprintf("[CHANGED] aws_sfn_state_machine %q: name mismatch (state: %q, live: %q)", arn, name, liveSM.Name))
		}
	}

	return drifts
}
