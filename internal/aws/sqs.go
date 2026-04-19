package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSQueue struct {
	URL  string
	Name string
	ARN  string
}

type sqsClient interface {
	ListQueues(ctx context.Context, params *sqs.ListQueuesInput, optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error)
	GetQueueAttributes(ctx context.Context, params *sqs.GetQueueAttributesInput, optFns ...func(*sqs.Options)) (*sqs.GetQueueAttributesOutput, error)
}

type SQSFetcher struct {
	client sqsClient
}

func NewSQSFetcher(client sqsClient) *SQSFetcher {
	return &SQSFetcher{client: client}
}

func (f *SQSFetcher) FetchAll(ctx context.Context) ([]SQSQueue, error) {
	out, err := f.client.ListQueues(ctx, &sqs.ListQueuesInput{})
	if err != nil {
		return nil, err
	}

	var queues []SQSQueue
	for _, url := range out.QueueUrls {
		attrs, err := f.client.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
			QueueUrl:       aws.String(url),
			AttributeNames: []string{"QueueArn"},
		})
		if err != nil {
			return nil, err
		}
		arn := attrs.Attributes["QueueArn"]
		name := extractQueueName(url)
		queues = append(queues, SQSQueue{URL: url, Name: name, ARN: arn})
	}
	return queues, nil
}

func extractQueueName(url string) string {
	for i := len(url) - 1; i >= 0; i-- {
		if url[i] == '/' {
			return url[i+1:]
		}
	}
	return url
}
