package aws_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	"github.com/stretchr/testify/assert"

	internalaws "driftctl-lite/internal/aws"
)

type mockELBClassicClient struct {
	output *elasticloadbalancing.DescribeLoadBalancersOutput
	err    error
}

func (m *mockELBClassicClient) DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancing.DescribeLoadBalancersInput, optFns ...func(*elasticloadbalancing.Options)) (*elasticloadbalancing.DescribeLoadBalancersOutput, error) {
	return m.output, m.err
}

func TestELBClassicFetchAll_ReturnsLoadBalancers(t *testing.T) {
	mock := &mockELBClassicClient{
		output: &elasticloadbalancing.DescribeLoadBalancersOutput{
			LoadBalancerDescriptions: []types.LoadBalancerDescription{
				{
					LoadBalancerName: aws.String("my-elb"),
					DNSName:          aws.String("my-elb.us-east-1.elb.amazonaws.com"),
					Scheme:           aws.String("internet-facing"),
				},
			},
		},
	}

	fetcher := internalaws.NewELBClassicFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "my-elb", results[0].Name)
	assert.Equal(t, "my-elb.us-east-1.elb.amazonaws.com", results[0].DNSName)
	assert.Equal(t, "internet-facing", results[0].Scheme)
	assert.Equal(t, "active", results[0].State)
}

func TestELBClassicFetchAll_Empty(t *testing.T) {
	mock := &mockELBClassicClient{
		output: &elasticloadbalancing.DescribeLoadBalancersOutput{
			LoadBalancerDescriptions: []types.LoadBalancerDescription{},
		},
	}

	fetcher := internalaws.NewELBClassicFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, results)
}
