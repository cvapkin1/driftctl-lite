package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
)

type mockSQSClient struct {
	listOut  *sqs.ListQueuesOutput
	attrsOut *sqs.GetQueueAttributesOutput
}

func (m *mockSQSClient) ListQueues(ctx context.Context, params *sqs.ListQueuesInput, optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error) {
	return m.listOut, nil
}

func (m *mockSQSClient) GetQueueAttributes(ctx context.Context, params *sqs.GetQueueAttributesInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error) {
	return m.attrsOut, nil
}

func TestSQSFetchAll_ReturnsQueues(t *testing.T) {
	mock := &mockSQSClient{
		listOut: &sqs.ListQueuesOutput{
			QueueUrls: []string{"https://sqs.us-east-1.amazonaws.com/123456789/my-queue"},
		},
		attrsOut: &sqs.GetQueueAttributesOutput{
			Attributes: map[string]string{"QueueArn": "arn:aws:sqs:us-east-1:123456789:my-queue"},
		},
	}
	fetcher := NewSQSFetcher(mock)
	queues, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, queues, 1)
	assert.Equal(t, "my-queue", queues[0].Name)
	assert.Equal(t, "arn:aws:sqs:us-east-1:123456789:my-queue", queues[0].ARN)
}

func TestSQSFetchAll_Empty(t *testing.T) {
	mock := &mockSQSClient{
		listOut:  &sqs.ListQueuesOutput{QueueUrls: []string{}},
		attrsOut: &sqs.GetQueueAttributesOutput{},
	}
	fetcher := NewSQSFetcher(mock)
	queues, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, queues)
}
