package drift

import (
	"fmt"

	"github.com/your-org/driftctl-lite/internal/aws"
)

type AppRunnerDriftResult struct {
	ServiceID string
	Issue     string
}

// DetectAppRunnerDrift compares Terraform state resources against live App Runner services.
// It flags services that are missing or in a non-running status.
func DetectAppRunnerDrift(stateResources []map[string]interface{}, live []aws.AppRunnerService) []AppRunnerDriftResult {
	var results []AppRunnerDriftResult

	liveByARN := make(map[string]aws.AppRunnerService, len(live))
	for _, svc := range live {
		liveByARN[svc.ARN] = svc
	}

	for _, res := range stateResources {
		arn, _ := res["arn"].(string)
		if arn == "" {
			continue
		}

		liveSvc, found := liveByARN[arn]
		if !found {
			results = append(results, AppRunnerDriftResult{
				ServiceID: arn,
				Issue:     "service not found in live AWS account",
			})
			continue
		}

		const runningStatus = "RUNNING"
		if liveSvc.Status != runningStatus {
			results = append(results, AppRunnerDriftResult{
				ServiceID: arn,
				Issue:     fmt.Sprintf("unexpected status: %s (expected %s)", liveSvc.Status, runningStatus),
			})
		}

		if stateURL, ok := res["service_url"].(string); ok && stateURL != "" {
			if liveSvc.URL != stateURL {
				results = append(results, AppRunnerDriftResult{
					ServiceID: arn,
					Issue:     fmt.Sprintf("service URL mismatch: state=%s live=%s", stateURL, liveSvc.URL),
				})
			}
		}
	}

	return results
}
