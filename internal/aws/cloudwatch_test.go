package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type mockCWClient struct {
	alarms []types.MetricAlarm
}

func (m *mockCWClient) DescribeAlarms(ctx context.Context, params *cloudwatch.DescribeAlarmsInput, optFns ...func(*cloudwatch.Options)) (*cloudwatch.DescribeAlarmsOutput, error) {
	return &cloudwatch.DescribeAlarmsOutput{MetricAlarms: m.alarms}, nil
}

func TestCloudWatchFetchAll_ReturnsAlarms(t *testing.T) {
	mock := &mockCWClient{
		alarms: []types.MetricAlarm{
			{AlarmName: aws.String("alarm-1"), AlarmArn: aws.String("arn:aws:cloudwatch::alarm-1"), StateValue: types.StateValueOk},
			{AlarmName: aws.String("alarm-2"), AlarmArn: aws.String("arn:aws:cloudwatch::alarm-2"), StateValue: types.StateValueAlarm},
		},
	}
	fetcher := NewCloudWatchFetcher(mock)
	alarms, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alarms) != 2 {
		t.Fatalf("expected 2 alarms, got %d", len(alarms))
	}
	if alarms[0].Name != "alarm-1" {
		t.Errorf("expected alarm-1, got %s", alarms[0].Name)
	}
	if alarms[1].State != "ALARM" {
		t.Errorf("expected ALARM state, got %s", alarms[1].State)
	}
}

func TestCloudWatchFetchAll_Empty(t *testing.T) {
	mock := &mockCWClient{alarms: []types.MetricAlarm{}}
	fetcher := NewCloudWatchFetcher(mock)
	alarms, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(alarms) != 0 {
		t.Errorf("expected 0 alarms, got %d", len(alarms))
	}
}
