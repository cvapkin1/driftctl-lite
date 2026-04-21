package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type mockCFNClient struct {
	output *cloudformation.DescribeStacksOutput
	err    error
}

func (m *mockCFNClient) DescribeStacks(_ context.Context, _ *cloudformation.DescribeStacksInput, _ ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error) {
	return m.output, m.err
}

func TestCloudFormationFetchAll_ReturnsStacks(t *testing.T) {
	mock := &mockCFNClient{
		output: &cloudformation.DescribeStacksOutput{
			Stacks: []types.Stack{
				{
					StackId:     aws.String("arn:aws:cloudformation:us-east-1:123456789012:stack/my-stack/abc123"),
					StackName:   aws.String("my-stack"),
					StackStatus: types.StackStatusCreateComplete,
				},
				{
					StackId:     aws.String("arn:aws:cloudformation:us-east-1:123456789012:stack/other-stack/def456"),
					StackName:   aws.String("other-stack"),
					StackStatus: types.StackStatusUpdateComplete,
				},
			},
		},
	}

	fetcher := NewCloudFormationFetcher(mock)
	stacks, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stacks) != 2 {
		t.Fatalf("expected 2 stacks, got %d", len(stacks))
	}
	if stacks[0].StackName != "my-stack" {
		t.Errorf("expected stack name 'my-stack', got '%s'", stacks[0].StackName)
	}
	if stacks[1].Status != "UPDATE_COMPLETE" {
		t.Errorf("expected status 'UPDATE_COMPLETE', got '%s'", stacks[1].Status)
	}
}

func TestCloudFormationFetchAll_Empty(t *testing.T) {
	mock := &mockCFNClient{
		output: &cloudformation.DescribeStacksOutput{
			Stacks: []types.Stack{},
		},
	}

	fetcher := NewCloudFormationFetcher(mock)
	stacks, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(stacks) != 0 {
		t.Errorf("expected 0 stacks, got %d", len(stacks))
	}
}
