package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSTopic struct {
	ARN  string
	Name string
	Tags map[string]string
}

type snsSvcAPI interface {
	ListTopics(ctx context.Context, params *sns.ListTopicsInput, optFns ...func(*sns.Options)) (*sns.ListTopicsOutput, error)
	ListTagsForResource(ctx context.Context, params *sns.ListTagsForResourceInput, optFns ...func(*sns.Options)) (*sns.ListTagsForResourceOutput, error)
}

type SNSFetcher struct {
	client snsSvcAPI
}

func NewSNSFetcher(client snsSvcAPI) *SNSFetcher {
	return &SNSFetcher{client: client}
}

func (f *SNSFetcher) FetchAll(ctx context.Context) ([]SNSTopic, error) {
	var topics []SNSTopic
	var nextToken *string

	for {
		out, err := f.client.ListTopics(ctx, &sns.ListTopicsInput{NextToken: nextToken})
		if err != nil {
			return nil, err
		}

		for _, t := range out.Topics {
			arn := aws.ToString(t.TopicArn)
			topic := SNSTopic{ARN: arn}

			tagOut, err := f.client.ListTagsForResource(ctx, &sns.ListTagsForResourceInput{
				ResourceArn: t.TopicArn,
			})
			if err == nil {
				topic.Tags = make(map[string]string, len(tagOut.Tags))
				for _, tag := range tagOut.Tags {
					topic.Tags[aws.ToString(tag.Key)] = aws.ToString(tag.Value)
				}
			}

			topics = append(topics, topic)
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return topics, nil
}
