package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	"github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/stretchr/testify/assert"
)

type mockConfigClient struct {
	output *configservice.DescribeConfigRulesOutput
	err    error
}

func (m *mockConfigClient) DescribeConfigRules(ctx context.Context, params *configservice.DescribeConfigRulesInput, optFns ...func(*configservice.Options)) (*configservice.DescribeConfigRulesOutput, error) {
	return m.output, m.err
}

func TestConfigFetchAll_ReturnsRules(t *testing.T) {
	mock := &mockConfigClient{
		output: &configservice.DescribeConfigRulesOutput{
			ConfigRules: []types.ConfigRule{
				{
					ConfigRuleName:  aws.String("rule-1"),
					ConfigRuleArn:   aws.String("arn:aws:config:us-east-1:123456789012:config-rule/rule-1"),
					ConfigRuleState: types.ConfigRuleStateActive,
					Source: &types.Source{
						Owner: types.OwnerAws,
					},
				},
			},
		},
	}

	fetcher := NewConfigFetcher(mock)
	rules, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, rules, 1)
	assert.Equal(t, "rule-1", rules[0].Name)
	assert.Equal(t, "ACTIVE", rules[0].State)
	assert.Equal(t, "AWS", rules[0].Source)
}

func TestConfigFetchAll_Empty(t *testing.T) {
	mock := &mockConfigClient{
		output: &configservice.DescribeConfigRulesOutput{
			ConfigRules: []types.ConfigRule{},
		},
	}

	fetcher := NewConfigFetcher(mock)
	rules, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, rules)
}
