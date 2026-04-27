package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type SSMParameter struct {
	Name    string
	ARN     string
	Type    string
	Version int64
}

type SSMClient interface {
	DescribeParameters(ctx context.Context, params *ssm.DescribeParametersInput, optFns ...func(*ssm.Options)) (*ssm.DescribeParametersOutput, error)
}

type SSMFetcher struct {
	client SSMClient
}

func NewSSMFetcher(client SSMClient) *SSMFetcher {
	return &SSMFetcher{client: client}
}

func (f *SSMFetcher) FetchAll(ctx context.Context) ([]SSMParameter, error) {
	var results []SSMParameter
	var nextToken *string

	for {
		out, err := f.client.DescribeParameters(ctx, &ssm.DescribeParametersInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, p := range out.Parameters {
			results = append(results, mapSSMParameter(p))
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return results, nil
}

func mapSSMParameter(p interface{ GetName() *string }) SSMParameter {
	// Use concrete type from SDK
	return SSMParameter{}
}

func mapSSMParameterMeta(name, arn, ptype string, version int64) SSMParameter {
	return SSMParameter{
		Name:    name,
		ARN:     arn,
		Type:    ptype,
		Version: version,
	}
}

func mapSSMFromSDK(name *string, arn *string, ptype string, version int64) SSMParameter {
	return SSMParameter{
		Name:    aws.ToString(name),
		ARN:     aws.ToString(arn),
		Type:    ptype,
		Version: version,
	}
}
