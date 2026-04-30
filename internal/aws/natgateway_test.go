package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func TestMapNATGateway_Fields(t *testing.T) {
	gw := types.NatGateway{
		NatGatewayId: aws.String("nat-0abc123"),
		State:        types.NatGatewayStateAvailable,
		VpcId:        aws.String("vpc-111"),
		SubnetId:     aws.String("subnet-222"),
		NatGatewayAddresses: []types.NatGatewayAddress{
			{PublicIp: aws.String("1.2.3.4")},
		},
	}

	result := mapNATGateway(gw)

	if result.ID != "nat-0abc123" {
		t.Errorf("expected ID nat-0abc123, got %s", result.ID)
	}
	if result.State != "available" {
		t.Errorf("expected state available, got %s", result.State)
	}
	if result.VPCID != "vpc-111" {
		t.Errorf("expected VPC vpc-111, got %s", result.VPCID)
	}
	if result.PublicIP != "1.2.3.4" {
		t.Errorf("expected public IP 1.2.3.4, got %s", result.PublicIP)
	}
}

func TestMapNATGateway_NoAddresses(t *testing.T) {
	gw := types.NatGateway{
		NatGatewayId:        aws.String("nat-empty"),
		State:               types.NatGatewayStatePending,
		NatGatewayAddresses: []types.NatGatewayAddress{},
	}

	result := mapNATGateway(gw)

	if result.PublicIP != "" {
		t.Errorf("expected empty PublicIP, got %s", result.PublicIP)
	}
	if result.State != "pending" {
		t.Errorf("expected state pending, got %s", result.State)
	}
}
