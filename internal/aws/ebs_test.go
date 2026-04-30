package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
)

type mockEBSClient struct {
	output *ec2.DescribeVolumesOutput
	err    error
}

func (m *mockEBSClient) DescribeVolumes(_ context.Context, _ *ec2.DescribeVolumesInput, _ ...func(*ec2.Options)) (*ec2.DescribeVolumesOutput, error) {
	return m.output, m.err
}

func TestEBSFetchAll_ReturnsVolumes(t *testing.T) {
	mock := &mockEBSClient{
		output: &ec2.DescribeVolumesOutput{
			Volumes: []types.Volume{
				{
					VolumeId:         aws.String("vol-abc123"),
					State:            types.VolumeStateAvailable,
					Size:             aws.Int32(100),
					VolumeType:       types.VolumeTypeGp3,
					AvailabilityZone: aws.String("us-east-1a"),
					Encrypted:        aws.Bool(true),
				},
			},
		},
	}
	fetcher := NewEBSFetcher(mock)
	result, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "vol-abc123", result[0].ID)
	assert.Equal(t, "available", result[0].State)
	assert.Equal(t, int32(100), result[0].Size)
	assert.Equal(t, "gp3", result[0].VolumeType)
	assert.Equal(t, "us-east-1a", result[0].AZ)
	assert.True(t, result[0].Encrypted)
}

func TestEBSFetchAll_Empty(t *testing.T) {
	mock := &mockEBSClient{
		output: &ec2.DescribeVolumesOutput{Volumes: []types.Volume{}},
	}
	fetcher := NewEBSFetcher(mock)
	result, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestMapEBSVolume_NilFields(t *testing.T) {
	v := types.Volume{
		VolumeId:  aws.String("vol-000"),
		State:     types.VolumeStateCreating,
		Encrypted: aws.Bool(false),
	}
	result := mapEBSVolume(v)
	assert.Equal(t, "vol-000", result.ID)
	assert.Equal(t, int32(0), result.Size)
	assert.Equal(t, "", result.AZ)
	assert.False(t, result.Encrypted)
}
