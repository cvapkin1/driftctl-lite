package drift

import (
	"fmt"

	"github.com/ekristen/driftctl-lite/internal/aws"
	"github.com/ekristen/driftctl-lite/internal/tfstate"
)

// DetectBudgetsDrift compares Terraform state budget resources against live AWS Budgets.
// It flags budgets that are missing, have a mismatched limit amount, or mismatched time unit.
func DetectBudgetsDrift(stateResources []tfstate.Resource, live []aws.BudgetResource) []DriftResult {
	var results []DriftResult

	liveIndex := make(map[string]aws.BudgetResource, len(live))
	for _, b := range live {
		liveIndex[b.Name] = b
	}

	for _, sr := range stateResources {
		if sr.Type != "aws_budgets_budget" {
			continue
		}

		name, _ := sr.Attributes["name"].(string)
		if name == "" {
			continue
		}

		liveB, found := liveIndex[name]
		if !found {
			results = append(results, DriftResult{
				ResourceType: sr.Type,
				ResourceID:   name,
				DriftType:    "missing",
				Details:      fmt.Sprintf("budget %q not found in AWS", name),
			})
			continue
		}

		stateLimitAmount, _ := sr.Attributes["limit_amount"].(string)
		if stateLimitAmount != "" && stateLimitAmount != liveB.LimitAmount {
			results = append(results, DriftResult{
				ResourceType: sr.Type,
				ResourceID:   name,
				DriftType:    "limit_amount_mismatch",
				Details:      fmt.Sprintf("budget %q limit_amount: state=%s live=%s", name, stateLimitAmount, liveB.LimitAmount),
			})
		}

		stateTimeUnit, _ := sr.Attributes["time_unit"].(string)
		if stateTimeUnit != "" && stateTimeUnit != liveB.TimeUnit {
			results = append(results, DriftResult{
				ResourceType: sr.Type,
				ResourceID:   name,
				DriftType:    "time_unit_mismatch",
				Details:      fmt.Sprintf("budget %q time_unit: state=%s live=%s", name, stateTimeUnit, liveB.TimeUnit),
			})
		}
	}

	return results
}
