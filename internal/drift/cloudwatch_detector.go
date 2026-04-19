package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

type CloudWatchDriftResult struct {
	AlarmName string
	Issue     string
}

func DetectCloudWatchDrift(resources []tfstate.Resource, live []aws.CloudWatchAlarm) []CloudWatchDriftResult {
	var results []CloudWatchDriftResult

	liveMap := make(map[string]aws.CloudWatchAlarm, len(live))
	for _, a := range live {
		liveMap[a.Name] = a
	}

	for _, r := range resources {
		if r.Type != "aws_cloudwatch_metric_alarm" {
			continue
		}
		name, _ := r.Attributes["alarm_name"].(string)
		if name == "" {
			continue
		}
		alarm, found := liveMap[name]
		if !found {
			results = append(results, CloudWatchDriftResult{
				AlarmName: name,
				Issue:     "alarm not found in AWS",
			})
			continue
		}
		if alarm.State == "ALARM" {
			results = append(results, CloudWatchDriftResult{
				AlarmName: name,
				Issue:     fmt.Sprintf("alarm is in ALARM state (may indicate drift)"),
			})
		}
	}
	return results
}
