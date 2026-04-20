package drift

import (
	"fmt"

	"github.com/evilmartians/driftctl-lite/internal/aws"
	"github.com/evilmartians/driftctl-lite/internal/tfstate"
)

func DetectELBv2Drift(stateResources []tfstate.Resource, live []aws.ELBv2Resource) []DriftResult {
	var results []DriftResult

	liveByARN := make(map[string]aws.ELBv2Resource, len(live))
	for _, lb := range live {
		liveByARN[lb.ARN] = lb
	}

	for _, sr := range stateResources {
		if sr.Type != "aws_lb" && sr.Type != "aws_alb" {
			continue
		}
		arn, _ := sr.Attributes["arn"].(string)
		if arn == "" {
			continue
		}
		lb, found := liveByARN[arn]
		if !found {
			results = append(results, DriftResult{
				ResourceType: sr.Type,
				ResourceID:   arn,
				DriftType:    "missing",
				Details:      "load balancer not found in AWS",
			})
			continue
		}
		if lb.State != "active" {
			results = append(results, DriftResult{
				ResourceType: sr.Type,
				ResourceID:   arn,
				DriftType:    "state_mismatch",
				Details:      fmt.Sprintf("expected active, got %s", lb.State),
			})
		}
	}
	return results
}
