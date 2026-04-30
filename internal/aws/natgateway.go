package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type NATGateway struct {
	ID        string
	State     string
	VPCID     string
	SubnetID  string
	PublicIP  string
}

type NATGatewayFetcher struct {
	client *ec2.Client
}

func NewNATGatewayFetcher(client *ec2.Client) *NATGatewayFetcher {
	return &NATGatewayFetcher{client: client}
}

func (f *NATGatewayFetcher) FetchAll(ctx context.Context) ([]NATGateway, error) {
	out, err := f.client.DescribeNatGateways(ctx, &ec2.DescribeNatGatewaysInput{})
	if err != nil {
		return nil, err
	}

	var gateways []NATGateway
	for _, gw := range out.NatGateways {
		gateways = append(gateways, mapNATGateway(gw))
	}
	return gateways, nil
}

func mapNATGateway(gw types.NatGateway) NATGateway {
	var publicIP string
	if len(gw.NatGatewayAddresses) > 0 && gw.NatGatewayAddresses[0].PublicIp != nil {
		publicIP = aws.ToString(gw.NatGatewayAddresses[0].PublicIp)
	}
	return NATGateway{
		ID:       aws.ToString(gw.NatGatewayId),
		State:    string(gw.State),
		VPCID:    aws.ToString(gw.VpcId),
		SubnetID: aws.ToString(gw.SubnetId),
		PublicIP: publicIP,
	}
}
