package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

// DetectGuardDutyDrift compares GuardDuty detectors from Terraform state
// against live AWS resources and returns a list of drift findings.
func DetectGuardDutyDrift(resources []tfstate.Resource, live []aws.GuardDutyDetector) []DriftResult {
	liveMap := make(map[string]aws.GuardDutyDetector, len(live))
	for _, d := range live {
		liveMap[d.ID] = d
	}

	var results []DriftResult

	for _, res := range resources {
		if res.Type != "aws_guardduty_detector" {
			continue
		}

		id, _ := res.Attributes["id"].(string)
		liveDetector, found := liveMap[id]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   id,
				ResourceType: res.Type,
				DriftType:    "missing",
				Details:      fmt.Sprintf("GuardDuty detector %s not found in AWS", id),
			})
			continue
		}

		if liveDetector.Status == "DISABLED" {
			results = append(results, DriftResult{
				ResourceID:   id,
				ResourceType: res.Type,
				DriftType:    "disabled",
				Details:      fmt.Sprintf("GuardDuty detector %s is DISABLED in AWS", id),
			})
			continue
		}

		stateEnable, _ := res.Attributes["enable"].(bool)
		if stateEnable && liveDetector.Status != "ENABLED" {
			results = append(results, DriftResult{
				ResourceID:   id,
				ResourceType: res.Type,
				DriftType:    "status_mismatch",
				Details:      fmt.Sprintf("GuardDuty detector %s: state=enabled, live=%s", id, liveDetector.Status),
			})
		}
	}

	return results
}
