package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

type LambdaDriftResult struct {
	ResourceID string
	DriftType  string
	Details    string
}

func DetectLambdaDrift(stateResources []tfstate.Resource, live []aws.LambdaFunction) []LambdaDriftResult {
	var results []LambdaDriftResult

	liveMap := make(map[string]aws.LambdaFunction)
	for _, fn := range live {
		liveMap[fn.Name] = fn
	}

	for _, res := range stateResources {
		if res.Type != "aws_lambda_function" {
			continue
		}
		name, _ := res.Attributes["function_name"].(string)
		liveFn, found := liveMap[name]
		if !found {
			results = append(results, LambdaDriftResult{
				ResourceID: name,
				DriftType:  "missing",
				Details:    "function not found in AWS",
			})
			continue
		}
		if liveFn.State != "Active" && liveFn.State != "active" {
			results = append(results, LambdaDriftResult{
				ResourceID: name,
				DriftType:  "state_mismatch",
				Details:    fmt.Sprintf("function state is %s", liveFn.State),
			})
		}
	}
	return results
}
