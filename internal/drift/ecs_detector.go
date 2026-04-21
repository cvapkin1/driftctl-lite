package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

const ecsInactiveStatus = "INACTIVE"

func DetectECSDrift(stateResources []tfstate.Resource, liveClusters []aws.ECSCluster) []DriftResult {
	var results []DriftResult

	liveMap := make(map[string]aws.ECSCluster)
	for _, c := range liveClusters {
		liveMap[c.ARN] = c
	}

	for _, res := range stateResources {
		if res.Type != "aws_ecs_cluster" {
			continue
		}

		arn, _ := res.Attributes["arn"].(string)
		if arn == "" {
			arn, _ = res.Attributes["id"].(string)
		}

		live, found := liveMap[arn]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   arn,
				ResourceType: res.Type,
				DriftType:    "missing",
				Details:      "ECS cluster not found in live AWS",
			})
			continue
		}

		if live.Status == ecsInactiveStatus {
			results = append(results, DriftResult{
				ResourceID:   arn,
				ResourceType: res.Type,
				DriftType:    "deleted",
				Details:      fmt.Sprintf("ECS cluster status is %s", live.Status),
			})
			continue
		}

		stateName, _ := res.Attributes["name"].(string)
		if stateName != "" && stateName != live.Name {
			results = append(results, DriftResult{
				ResourceID:   arn,
				ResourceType: res.Type,
				DriftType:    "modified",
				Details:      fmt.Sprintf("name mismatch: state=%s live=%s", stateName, live.Name),
			})
		}
	}

	return results
}
