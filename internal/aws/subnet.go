package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type SubnetResource struct {
	ID               string
	VPCID            string
	CIDRBlock        string
	AvailabilityZone string
	State            string
	MapPublicIPOnLaunch bool
}

type subnetEC2Client interface {
	DescribeSubnets(ctx context.Context, params *ec2.DescribeSubnetsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeSubnetsOutput, error)
}

type SubnetFetcher struct {
	client subnetEC2Client
}

func NewSubnetFetcher(client subnetEC2Client) *SubnetFetcher {
	return &SubnetFetcher{client: client}
}

func (f *SubnetFetcher) FetchAll(ctx context.Context) ([]SubnetResource, error) {
	output, err := f.client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{})
	if err != nil {
		return nil, err
	}

	resources := make([]SubnetResource, 0, len(output.Subnets))
	for _, s := range output.Subnets {
		resources = append(resources, mapSubnet(s))
	}
	return resources, nil
}

func mapSubnet(s types.Subnet) SubnetResource {
	return SubnetResource{
		ID:               aws.ToString(s.SubnetId),
		VPCID:            aws.ToString(s.VpcId),
		CIDRBlock:        aws.ToString(s.CidrBlock),
		AvailabilityZone: aws.ToString(s.AvailabilityZone),
		State:            string(s.State),
		MapPublicIPOnLaunch: aws.ToBool(s.MapPublicIpOnLaunch),
	}
}
