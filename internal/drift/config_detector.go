package drift

import (
	"fmt"

	"github.com/edobry/driftctl-lite/internal/aws"
)

type ConfigDriftResult struct {
	ResourceID string
	DriftType  string
	Details    string
}

func DetectConfigDrift(stateResources []map[string]interface{}, liveRules []aws.ConfigRule) []ConfigDriftResult {
	var results []ConfigDriftResult

	liveMap := make(map[string]aws.ConfigRule)
	for _, r := range liveRules {
		liveMap[r.Name] = r
	}

	for _, res := range stateResources {
		attrs, ok := res["attributes"].(map[string]interface{})
		if !ok {
			continue
		}

		name, _ := attrs["name"].(string)
		if name == "" {
			continue
		}

		live, found := liveMap[name]
		if !found {
			results = append(results, ConfigDriftResult{
				ResourceID: name,
				DriftType:  "missing",
				Details:    fmt.Sprintf("Config rule %q not found in AWS", name),
			})
			continue
		}

		if live.State == "DELETING" || live.State == "DELETING_RESULTS" {
			results = append(results, ConfigDriftResult{
				ResourceID: name,
				DriftType:  "deleted",
				Details:    fmt.Sprintf("Config rule %q is in state %q", name, live.State),
			})
			continue
		}

		if stateARN, ok := attrs["arn"].(string); ok && stateARN != "" && stateARN != live.ARN {
			results = append(results, ConfigDriftResult{
				ResourceID: name,
				DriftType:  "arn_mismatch",
				Details:    fmt.Sprintf("Config rule %q ARN mismatch: state=%q live=%q", name, stateARN, live.ARN),
			})
		}
	}

	return results
}
