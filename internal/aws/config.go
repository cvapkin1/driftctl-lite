package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
)

type ConfigRule struct {
	Name      string
	ARN       string
	State     string
	Source    string
}

type ConfigServiceClient interface {
	DescribeConfigRules(ctx context.Context, params *configservice.DescribeConfigRulesInput, optFns ...func(*configservice.Options)) (*configservice.DescribeConfigRulesOutput, error)
}

type ConfigFetcher struct {
	client ConfigServiceClient
}

func NewConfigFetcher(client ConfigServiceClient) *ConfigFetcher {
	return &ConfigFetcher{client: client}
}

func (f *ConfigFetcher) FetchAll(ctx context.Context) ([]ConfigRule, error) {
	var rules []ConfigRule
	var nextToken *string

	for {
		out, err := f.client.DescribeConfigRules(ctx, &configservice.DescribeConfigRulesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range out.ConfigRules {
			rules = append(rules, mapConfigRule(r))
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return rules, nil
}

func mapConfigRule(r types.ConfigRule) ConfigRule {
	name := aws.ToString(r.ConfigRuleName)
	arn := aws.ToString(r.ConfigRuleArn)
	state := string(r.ConfigRuleState)
	source := ""
	if r.Source != nil {
		source = aws.ToString(r.Source.Owner)
	}
	return ConfigRule{
		Name:   name,
		ARN:    arn,
		State:  state,
		Source: source,
	}
}
