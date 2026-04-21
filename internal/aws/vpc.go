package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type VPCResource struct {
	ID        string
	CIDR      string
	State     string
	IsDefault bool
	Tags      map[string]string
}

type vpcEC2Client interface {
	DescribeVpcs(ctx context.Context, params *ec2.DescribeVpcsInput, optFns ...func(*ec2.Options)) (*ec2.DescribeVpcsOutput, error)
}

type VPCFetcher struct {
	client vpcEC2Client
}

func NewVPCFetcher(client vpcEC2Client) *VPCFetcher {
	return &VPCFetcher{client: client}
}

func (f *VPCFetcher) FetchAll(ctx context.Context) ([]VPCResource, error) {
	out, err := f.client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, err
	}

	results := make([]VPCResource, 0, len(out.Vpcs))
	for _, v := range out.Vpcs {
		results = append(results, mapVPC(v))
	}
	return results, nil
}

func mapVPC(v types.Vpc) VPCResource {
	tags := make(map[string]string)
	for _, t := range v.Tags {
		if t.Key != nil && t.Value != nil {
			tags[*t.Key] = *t.Value
		}
	}
	return VPCResource{
		ID:        aws.ToString(v.VpcId),
		CIDR:      aws.ToString(v.CidrBlock),
		State:     string(v.State),
		IsDefault: aws.ToBool(v.IsDefault),
		Tags:      tags,
	}
}
