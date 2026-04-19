package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
)

type mockSNSClient struct {
	topics []types.Topic
	tags   map[string][]types.Tag
}

func (m *mockSNSClient) ListTopics(ctx context.Context, params *sns.ListTopicsInput, optFns ...func(*sns.Options)) (*sns.ListTopicsOutput, error) {
	return &sns.ListTopicsOutput{Topics: m.topics}, nil
}

func (m *mockSNSClient) ListTagsForResource(ctx context.Context, params *sns.ListTagsForResourceInput, optFns ...func(*sns.Options)) (*sns.ListTagsForResourceOutput, error) {
	arn := aws.ToString(params.ResourceArn)
	return &sns.ListTagsForResourceOutput{Tags: m.tags[arn]}, nil
}

func TestSNSFetchAll_ReturnsTopics(t *testing.T) {
	mock := &mockSNSClient{
		topics: []types.Topic{
			{TopicArn: aws.String("arn:aws:sns:us-east-1:123456789012:my-topic")},
		},
		tags: map[string][]types.Tag{
			"arn:aws:sns:us-east-1:123456789012:my-topic": {
				{Key: aws.String("Env"), Value: aws.String("prod")},
			},
		},
	}

	fetcher := NewSNSFetcher(mock)
	topics, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(topics) != 1 {
		t.Fatalf("expected 1 topic, got %d", len(topics))
	}
	if topics[0].ARN != "arn:aws:sns:us-east-1:123456789012:my-topic" {
		t.Errorf("unexpected ARN: %s", topics[0].ARN)
	}
	if topics[0].Tags["Env"] != "prod" {
		t.Errorf("expected tag Env=prod")
	}
}

func TestSNSFetchAll_Empty(t *testing.T) {
	mock := &mockSNSClient{}
	fetcher := NewSNSFetcher(mock)
	topics, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(topics) != 0 {
		t.Errorf("expected 0 topics, got %d", len(topics))
	}
}
