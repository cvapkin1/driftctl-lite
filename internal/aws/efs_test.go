package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
	"github.com/stretchr/testify/assert"
)

type mockEFSClient struct {
	output *efs.DescribeFileSystemsOutput
	err    error
}

func (m *mockEFSClient) DescribeFileSystems(_ context.Context, _ *efs.DescribeFileSystemsInput, _ ...func(*efs.Options)) (*efs.DescribeFileSystemsOutput, error) {
	return m.output, m.err
}

func TestEFSFetchAll_ReturnsFileSystems(t *testing.T) {
	mock := &mockEFSClient{
		output: &efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{
				{
					FileSystemId:   aws.String("fs-abc123"),
					Name:           aws.String("my-efs"),
					LifeCycleState: efstypes.LifeCycleStateAvailable,
					Encrypted:      aws.Bool(true),
					ThroughputMode: efstypes.ThroughputModeBursting,
				},
			},
		},
	}
	fetcher := &EFSFetcher{client: mock}
	results, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "fs-abc123", results[0].ID)
	assert.Equal(t, "my-efs", results[0].Name)
	assert.Equal(t, "available", results[0].LifeCycleState)
	assert.True(t, results[0].Encrypted)
	assert.Equal(t, "bursting", results[0].ThroughputMode)
}

func TestEFSFetchAll_Empty(t *testing.T) {
	mock := &mockEFSClient{
		output: &efs.DescribeFileSystemsOutput{
			FileSystems: []efstypes.FileSystemDescription{},
		},
	}
	fetcher := &EFSFetcher{client: mock}
	results, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, results)
}
