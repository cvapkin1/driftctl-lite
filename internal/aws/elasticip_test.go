package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockElasticIPClient struct {
	output *ec2.DescribeAddressesOutput
	err    error
}

func (m *mockElasticIPClient) DescribeAddresses(ctx context.Context, params *ec2.DescribeAddressesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeAddressesOutput, error) {
	return m.output, m.err
}

func TestElasticIPFetchAll_ReturnsAddresses(t *testing.T) {
	mock := &mockElasticIPClient{
		output: &ec2.DescribeAddressesOutput{
			Addresses: []types.Address{
				{
					AllocationId:  aws.String("eipalloc-abc123"),
					PublicIp:      aws.String("1.2.3.4"),
					InstanceId:    aws.String("i-0abc123"),
					Domain:        types.DomainTypeVpc,
					AssociationId: aws.String("eipassoc-xyz"),
				},
			},
		},
	}
	fetcher := NewElasticIPFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	require.NoError(t, err)
	require.Len(t, results, 1)
	assert.Equal(t, "eipalloc-abc123", results[0].AllocationID)
	assert.Equal(t, "1.2.3.4", results[0].PublicIP)
	assert.Equal(t, "i-0abc123", results[0].InstanceID)
	assert.Equal(t, "vpc", results[0].Domain)
	assert.True(t, results[0].Associated)
}

func TestElasticIPFetchAll_Empty(t *testing.T) {
	mock := &mockElasticIPClient{
		output: &ec2.DescribeAddressesOutput{},
	}
	fetcher := NewElasticIPFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	require.NoError(t, err)
	assert.Empty(t, results)
}

func TestMapElasticIP_Unassociated(t *testing.T) {
	addr := types.Address{
		AllocationId: aws.String("eipalloc-unassoc"),
		PublicIp:     aws.String("5.6.7.8"),
		Domain:       types.DomainTypeVpc,
	}
	res := mapElasticIP(addr)
	assert.Equal(t, "eipalloc-unassoc", res.AllocationID)
	assert.False(t, res.Associated)
	assert.Empty(t, res.InstanceID)
}
