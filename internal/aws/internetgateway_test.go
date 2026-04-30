package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
)

type mockIGWClient struct {
	output *ec2.DescribeInternetGatewaysOutput
	err    error
}

func (m *mockIGWClient) DescribeInternetGateways(_ context.Context, _ *ec2.DescribeInternetGatewaysInput, _ ...func(*ec2.Options)) (*ec2.DescribeInternetGatewaysOutput, error) {
	return m.output, m.err
}

func TestInternetGatewayFetchAll_ReturnsGateways(t *testing.T) {
	mock := &mockIGWClient{
		output: &ec2.DescribeInternetGatewaysOutput{
			InternetGateways: []types.InternetGateway{
				{
					InternetGatewayId: aws.String("igw-abc123"),
					OwnerId:           aws.String("123456789012"),
					Attachments: []types.InternetGatewayAttachment{
						{VpcId: aws.String("vpc-111"), State: types.AttachmentStatusAttached},
					},
				},
			},
		},
	}
	fetcher := NewInternetGatewayFetcher(mock)
	result, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "igw-abc123", result[0].ID)
	assert.Equal(t, "vpc-111", result[0].VPCID)
	assert.Equal(t, "attached", result[0].State)
}

func TestInternetGatewayFetchAll_Empty(t *testing.T) {
	mock := &mockIGWClient{
		output: &ec2.DescribeInternetGatewaysOutput{},
	}
	fetcher := NewInternetGatewayFetcher(mock)
	result, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestMapInternetGateway_Detached(t *testing.T) {
	igw := types.InternetGateway{
		InternetGatewayId: aws.String("igw-detached"),
		OwnerId:           aws.String("999"),
		Attachments:       []types.InternetGatewayAttachment{},
	}
	result := mapInternetGateway(igw)
	assert.Equal(t, "igw-detached", result.ID)
	assert.Equal(t, "", result.VPCID)
	assert.Equal(t, "detached", result.State)
}
