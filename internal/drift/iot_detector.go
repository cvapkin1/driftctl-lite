package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

type IoTDriftResult struct {
	ResourceID string
	Status     string
	Details    string
}

func DetectIoTDrift(stateResources []tfstate.Resource, liveThings []aws.IoTThing) []IoTDriftResult {
	var results []IoTDriftResult

	liveMap := make(map[string]aws.IoTThing)
	for _, t := range liveThings {
		liveMap[t.ThingName] = t
	}

	for _, res := range stateResources {
		if res.Type != "aws_iot_thing" {
			continue
		}

		name, _ := res.Attributes["name"].(string)
		expectedARN, _ := res.Attributes["arn"].(string)

		live, found := liveMap[name]
		if !found {
			results = append(results, IoTDriftResult{
				ResourceID: name,
				Status:     "missing",
				Details:    fmt.Sprintf("IoT thing %q not found in AWS", name),
			})
			continue
		}

		if expectedARN != "" && live.ARN != expectedARN {
			results = append(results, IoTDriftResult{
				ResourceID: name,
				Status:     "drifted",
				Details:    fmt.Sprintf("ARN mismatch: state=%q live=%q", expectedARN, live.ARN),
			})
			continue
		}

		results = append(results, IoTDriftResult{
			ResourceID: name,
			Status:     "ok",
			Details:    "",
		})
	}

	return results
}
