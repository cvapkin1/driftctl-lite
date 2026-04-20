package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stretchr/testify/assert"
)

type mockEKSClient struct {
	listOut     *eks.ListClustersOutput
	describeOut *eks.DescribeClusterOutput
	err         error
}

func (m *mockEKSClient) ListClusters(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error) {
	return m.listOut, m.err
}

func (m *mockEKSClient) DescribeCluster(ctx context.Context, params *eks.DescribeClusterInput, optFns ...func(*eks.Options)) (*eks.DescribeClusterOutput, error) {
	return m.describeOut, m.err
}

func TestEKSFetchAll_ReturnsClusters(t *testing.T) {
	mock := &mockEKSClient{
		listOut: &eks.ListClustersOutput{
			Clusters: []string{"my-cluster"},
		},
		describeOut: &eks.DescribeClusterOutput{
			Cluster: &eksTypes.Cluster{
				Name:    aws.String("my-cluster"),
				Arn:     aws.String("arn:aws:eks:us-east-1:123456789012:cluster/my-cluster"),
				Status:  eksTypes.ClusterStatusActive,
				Version: aws.String("1.29"),
			},
		},
	}

	fetcher := NewEKSFetcher(mock)
	clusters, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, clusters, 1)
	assert.Equal(t, "my-cluster", clusters[0].Name)
	assert.Equal(t, "ACTIVE", clusters[0].Status)
	assert.Equal(t, "1.29", clusters[0].Version)
}

func TestEKSFetchAll_Empty(t *testing.T) {
	mock := &mockEKSClient{
		listOut: &eks.ListClustersOutput{
			Clusters: []string{},
		},
	}

	fetcher := NewEKSFetcher(mock)
	clusters, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, clusters)
}
