package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"
)

type ELBv2Resource struct {
	ARN    string
	Name   string
	Type   string
	State  string
	DNS    string
}

type elbv2Client interface {
	DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancingv2.DescribeLoadBalancersInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeLoadBalancersOutput, error)
}

type ELBv2Fetcher struct {
	client elbv2Client
}

func NewELBv2Fetcher(client elbv2Client) *ELBv2Fetcher {
	return &ELBv2Fetcher{client: client}
}

func (f *ELBv2Fetcher) FetchAll(ctx context.Context) ([]ELBv2Resource, error) {
	var resources []ELBv2Resource
	var marker *string

	for {
		out, err := f.client.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{
			Marker: marker,
		})
		if err != nil {
			return nil, err
		}
		for _, lb := range out.LoadBalancers {
			resources = append(resources, mapLB(lb))
		}
		if out.NextMarker == nil {
			break
		}
		marker = out.NextMarker
	}
	return resources, nil
}

func mapLB(lb types.LoadBalancer) ELBv2Resource {
	return ELBv2Resource{
		ARN:   aws.ToString(lb.LoadBalancerArn),
		Name:  aws.ToString(lb.LoadBalancerName),
		Type:  string(lb.Type),
		State: string(lb.State.Code),
		DNS:   aws.ToString(lb.DNSName),
	}
}
