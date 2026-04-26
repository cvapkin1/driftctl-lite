package aws_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kafka/types"
	"github.com/stretchr/testify/assert"
)

type mockMSKClient struct {
	clusters []types.ClusterInfo
	err      error
}

func (m *mockMSKClient) ListClusters(ctx context.Context, params *kafka.ListClustersInput, optFns ...func(*kafka.Options)) (*kafka.ListClustersOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &kafka.ListClustersOutput{
		ClusterInfoList: m.clusters,
	}, nil
}

func TestMSKFetchAll_ReturnsClusters(t *testing.T) {
	mock := &mockMSKClient{
		clusters: []types.ClusterInfo{
			{
				ClusterArn:  aws.String("arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abc-123"),
				ClusterName: aws.String("my-cluster"),
				State:       types.ClusterStateActive,
				CurrentBrokerSoftwareInfo: &types.BrokerSoftwareInfo{
					KafkaVersion: aws.String("2.8.1"),
				},
			},
			{
				ClusterArn:  aws.String("arn:aws:kafka:us-east-1:123456789012:cluster/other-cluster/def-456"),
				ClusterName: aws.String("other-cluster"),
				State:       types.ClusterStateCreating,
				CurrentBrokerSoftwareInfo: &types.BrokerSoftwareInfo{
					KafkaVersion: aws.String("3.3.1"),
				},
			},
		},
	}

	fetcher := newMSKFetcherWithClient(mock)
	results, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, results, 2)

	assert.Equal(t, "arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abc-123", results[0].ID)
	assert.Equal(t, "my-cluster", results[0].Name)
	assert.Equal(t, "ACTIVE", results[0].State)
	assert.Equal(t, "2.8.1", results[0].KafkaVersion)

	assert.Equal(t, "arn:aws:kafka:us-east-1:123456789012:cluster/other-cluster/def-456", results[1].ID)
	assert.Equal(t, "other-cluster", results[1].Name)
	assert.Equal(t, "CREATING", results[1].State)
	assert.Equal(t, "3.3.1", results[1].KafkaVersion)
}

func TestMSKFetchAll_Empty(t *testing.T) {
	mock := &mockMSKClient{
		clusters: []types.ClusterInfo{},
	}

	fetcher := newMSKFetcherWithClient(mock)
	results, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, results)
}

func TestMSKFetchAll_NilKafkaVersion(t *testing.T) {
	mock := &mockMSKClient{
		clusters: []types.ClusterInfo{
			{
				ClusterArn:                aws.String("arn:aws:kafka:us-east-1:123456789012:cluster/bare-cluster/ghi-789"),
				ClusterName:               aws.String("bare-cluster"),
				State:                     types.ClusterStateActive,
				CurrentBrokerSoftwareInfo: nil,
			},
		},
	}

	fetcher := newMSKFetcherWithClient(mock)
	results, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "", results[0].KafkaVersion)
}
