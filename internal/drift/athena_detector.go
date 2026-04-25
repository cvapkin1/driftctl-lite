package drift

import (
	"fmt"

	"github.com/edasaki/driftctl-lite/internal/aws"
)

type AthenaWorkgroupState struct {
	Name   string
	State  string
	Engine string
}

func DetectAthenaDrift(stateWorkgroups []AthenaWorkgroupState, liveWorkgroups []aws.AthenaWorkgroup) []DriftResult {
	var results []DriftResult

	liveMap := make(map[string]aws.AthenaWorkgroup)
	for _, wg := range liveWorkgroups {
		liveMap[wg.Name] = wg
	}

	for _, stateWG := range stateWorkgroups {
		live, found := liveMap[stateWG.Name]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   stateWG.Name,
				ResourceType: "aws_athena_workgroup",
				DriftType:    "missing",
				Details:      fmt.Sprintf("workgroup %q not found in AWS", stateWG.Name),
			})
			continue
		}

		if live.State == "DISABLED" {
			results = append(results, DriftResult{
				ResourceID:   stateWG.Name,
				ResourceType: "aws_athena_workgroup",
				DriftType:    "disabled",
				Details:      fmt.Sprintf("workgroup %q is DISABLED in AWS", stateWG.Name),
			})
			continue
		}

		if stateWG.Engine != "" && live.Engine != stateWG.Engine {
			results = append(results, DriftResult{
				ResourceID:   stateWG.Name,
				ResourceType: "aws_athena_workgroup",
				DriftType:    "engine_mismatch",
				Details:      fmt.Sprintf("workgroup %q engine: state=%q live=%q", stateWG.Name, stateWG.Engine, live.Engine),
			})
		}
	}

	return results
}
