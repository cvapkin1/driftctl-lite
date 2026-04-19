package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type CloudWatchAlarm struct {
	Name  string
	State string
	ARN   string
}

type CloudWatchClient interface {
	DescribeAlarms(ctx context.Context, params *cloudwatch.DescribeAlarmsInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.DescribeAlarmsOutput, error)
}

type CloudWatchFetcher struct {
	client CloudWatchClient
}

func NewCloudWatchFetcher(client CloudWatchClient) *CloudWatchFetcher {
	return &CloudWatchFetcher{client: client}
}

func (f *CloudWatchFetcher) FetchAll(ctx context.Context) ([]CloudWatchAlarm, error) {
	var alarms []CloudWatchAlarm
	var nextToken *string

	for {
		out, err := f.client.DescribeAlarms(ctx, &cloudwatch.DescribeAlarmsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, a := range out.MetricAlarms {
			alarms = append(alarms, mapAlarm(a))
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return alarms, nil
}

func mapAlarm(a types.MetricAlarm) CloudWatchAlarm {
	return CloudWatchAlarm{
		Name:  aws.ToString(a.AlarmName),
		State: string(a.StateValue),
		ARN:   aws.ToString(a.AlarmArn),
	}
}
