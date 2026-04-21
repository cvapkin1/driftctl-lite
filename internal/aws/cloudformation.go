package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
)

type CloudFormationStack struct {
	StackID   string
	StackName string
	Status    string
}

type cloudFormationClient interface {
	DescribeStacks(ctx context.Context, params *cloudformation.DescribeStacksInput, optFns ...func(*cloudformation.Options)) (*cloudformation.DescribeStacksOutput, error)
}

type CloudFormationFetcher struct {
	client cloudFormationClient
}

func NewCloudFormationFetcher(client cloudFormationClient) *CloudFormationFetcher {
	return &CloudFormationFetcher{client: client}
}

func (f *CloudFormationFetcher) FetchAll(ctx context.Context) ([]CloudFormationStack, error) {
	var stacks []CloudFormationStack
	var nextToken *string

	for {
		out, err := f.client.DescribeStacks(ctx, &cloudformation.DescribeStacksInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, s := range out.Stacks {
			stacks = append(stacks, mapStack(s))
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return stacks, nil
}

func mapStack(s types.Stack) CloudFormationStack {
	return CloudFormationStack{
		StackID:   aws.ToString(s.StackId),
		StackName: aws.ToString(s.StackName),
		Status:    string(s.StackStatus),
	}
}
