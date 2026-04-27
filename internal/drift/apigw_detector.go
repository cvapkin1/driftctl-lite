package drift

import (
	"fmt"

	"github.com/example/driftctl-lite/internal/aws"
)

type APIGatewayDriftResult struct {
	ResourceID string
	Status     string
	Details    string
}

// DetectAPIGatewayDrift compares Terraform state resources against live API Gateway REST APIs.
// It reports missing or ARN-mismatched APIs as drift.
func DetectAPIGatewayDrift(stateResources []map[string]interface{}, live []aws.APIGatewayResource) []APIGatewayDriftResult {
	liveByID := make(map[string]aws.APIGatewayResource, len(live))
	for _, api := range live {
		liveByID[api.ID] = api
	}

	var results []APIGatewayDriftResult

	for _, res := range stateResources {
		id, _ := res["id"].(string)
		expectedARN, _ := res["arn"].(string)

		liveAPI, found := liveByID[id]
		if !found {
			results = append(results, APIGatewayDriftResult{
				ResourceID: id,
				Status:     "missing",
				Details:    fmt.Sprintf("REST API %q not found in AWS", id),
			})
			continue
		}

		if expectedARN != "" && liveAPI.ARN != expectedARN {
			results = append(results, APIGatewayDriftResult{
				ResourceID: id,
				Status:     "drifted",
				Details:    fmt.Sprintf("ARN mismatch: state=%q live=%q", expectedARN, liveAPI.ARN),
			})
			continue
		}

		results = append(results, APIGatewayDriftResult{
			ResourceID: id,
			Status:     "ok",
			Details:    "",
		})
	}

	return results
}
