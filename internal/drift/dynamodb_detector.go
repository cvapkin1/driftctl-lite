package drift

import (
	"fmt"

	"github.com/acme/driftctl-lite/internal/aws"
)

type DynamoDBDriftResult struct {
	TableName string
	DriftType string
	Details   string
}

func DetectDynamoDBDrift(stateResources []map[string]interface{}, live []aws.DynamoDBTable) []DynamoDBDriftResult {
	var results []DynamoDBDriftResult

	liveMap := make(map[string]aws.DynamoDBTable)
	for _, t := range live {
		liveMap[t.Name] = t
	}

	for _, res := range stateResources {
		name, _ := res["name"].(string)
		arn, _ := res["arn"].(string)

		liveTable, found := liveMap[name]
		if !found {
			results = append(results, DynamoDBDriftResult{
				TableName: name,
				DriftType: "missing",
				Details:   "table not found in live AWS account",
			})
			continue
		}

		if liveTable.Status == "DELETING" || liveTable.Status == "" {
			results = append(results, DynamoDBDriftResult{
				TableName: name,
				DriftType: "deleted",
				Details:   fmt.Sprintf("table status is %q", liveTable.Status),
			})
			continue
		}

		if arn != "" && liveTable.ARN != arn {
			results = append(results, DynamoDBDriftResult{
				TableName: name,
				DriftType: "arn_mismatch",
				Details:   fmt.Sprintf("state ARN %q != live ARN %q", arn, liveTable.ARN),
			})
		}
	}

	return results
}
