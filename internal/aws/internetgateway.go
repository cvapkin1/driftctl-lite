package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type InternetGateway struct {
	ID        string
	VPCID     string
	State     string
	OwnerID   string
}

type igwEC2Client interface {
	DescribeInternetGateways(ctx context.Context, params *ec2.DescribeInternetGatewaysInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInternetGatewaysOutput, error)
}

type InternetGatewayFetcher struct {
	client igwEC2Client
}

func NewInternetGatewayFetcher(client igwEC2Client) *InternetGatewayFetcher {
	return &InternetGatewayFetcher{client: client}
}

func (f *InternetGatewayFetcher) FetchAll(ctx context.Context) ([]InternetGateway, error) {
	out, err := f.client.DescribeInternetGateways(ctx, &ec2.DescribeInternetGatewaysInput{})
	if err != nil {
		return nil, err
	}
	var gateways []InternetGateway
	for _, igw := range out.InternetGateways {
		gateways = append(gateways, mapInternetGateway(igw))
	}
	return gateways, nil
}

func mapInternetGateway(igw types.InternetGateway) InternetGateway {
	vpcID := ""
	state := "detached"
	if len(igw.Attachments) > 0 {
		if igw.Attachments[0].VpcId != nil {
			vpcID = aws.ToString(igw.Attachments[0].VpcId)
		}
		state = string(igw.Attachments[0].State)
	}
	return InternetGateway{
		ID:      aws.ToString(igw.InternetGatewayId),
		VPCID:   vpcID,
		State:   state,
		OwnerID: aws.ToString(igw.OwnerId),
	}
}
