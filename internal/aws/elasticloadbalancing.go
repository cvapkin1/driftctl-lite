package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
)

type ELBClassicResource struct {
	Name      string
	DNSName   string
	Scheme    string
	State     string
}

type elbClassicClient interface {
	DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancing.DescribeLoadBalancersInput, optFns ...func(*elasticloadbalancing.Options)) (*elasticloadbalancing.DescribeLoadBalancersOutput, error)
}

type ELBClassicFetcher struct {
	client elbClassicClient
}

func NewELBClassicFetcher(client elbClassicClient) *ELBClassicFetcher {
	return &ELBClassicFetcher{client: client}
}

func (f *ELBClassicFetcher) FetchAll(ctx context.Context) ([]ELBClassicResource, error) {
	out, err := f.client.DescribeLoadBalancers(ctx, &elasticloadbalancing.DescribeLoadBalancersInput{})
	if err != nil {
		return nil, err
	}

	var resources []ELBClassicResource
	for _, lb := range out.LoadBalancerDescriptions {
		resources = append(resources, mapELBClassic(lb))
	}
	return resources, nil
}

func mapELBClassic(lb types.LoadBalancerDescription) ELBClassicResource {
	name := aws.ToString(lb.LoadBalancerName)
	dns := aws.ToString(lb.DNSName)
	scheme := aws.ToString(lb.Scheme)

	state := "active"
	if name == "" {
		state = "unknown"
	}

	return ELBClassicResource{
		Name:    name,
		DNSName: dns,
		Scheme:  scheme,
		State:   state,
	}
}
