package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
	"github.com/aws/aws-sdk-go-v2/service/opensearch/types"
	"github.com/stretchr/testify/assert"
)

// mockOpenSearchClient implements a minimal OpenSearch client for testing.
type mockOpenSearchClient struct {
	domains []types.DomainInfo
	statuses map[string]*types.DomainStatus
}

func (m *mockOpenSearchClient) ListDomainNames(ctx context.Context, params *opensearch.ListDomainNamesInput, optFns ...func(*opensearch.Options)) (*opensearch.ListDomainNamesOutput, error) {
	return &opensearch.ListDomainNamesOutput{DomainNames: m.domains}, nil
}

func (m *mockOpenSearchClient) DescribeDomains(ctx context.Context, params *opensearch.DescribeDomainsInput, optFns ...func(*opensearch.Options)) (*opensearch.DescribeDomainsOutput, error) {
	var statuses []types.DomainStatus
	for _, name := range params.DomainNames {
		if s, ok := m.statuses[name]; ok {
			statuses = append(statuses, *s)
		}
	}
	return &opensearch.DescribeDomainsOutput{DomainStatusList: statuses}, nil
}

func TestOpenSearchFetchAll_ReturnsDomains(t *testing.T) {
	mock := &mockOpenSearchClient{
		domains: []types.DomainInfo{
			{DomainName: aws.String("my-domain")},
		},
		statuses: map[string]*types.DomainStatus{
			"my-domain": {
				DomainName: aws.String("my-domain"),
				ARN:        aws.String("arn:aws:es:us-east-1:123456789012:domain/my-domain"),
				EngineVersion: aws.String("OpenSearch_2.3"),
				Deleted: aws.Bool(false),
				ClusterConfig: &types.ClusterConfig{
					InstanceType: types.OpenSearchPartitionInstanceTypeM5LargeSearch,
				},
			},
		},
	}

	fetcher := &OpenSearchFetcher{client: mock}
	resources, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, resources, 1)
	assert.Equal(t, "my-domain", resources[0].ID)
	assert.Equal(t, "arn:aws:es:us-east-1:123456789012:domain/my-domain", resources[0].ARN)
	assert.Equal(t, "OpenSearch_2.3", resources[0].EngineVersion)
	assert.Equal(t, string(types.OpenSearchPartitionInstanceTypeM5LargeSearch), resources[0].InstanceType)
	assert.False(t, resources[0].Deleted)
}

func TestOpenSearchFetchAll_Empty(t *testing.T) {
	mock := &mockOpenSearchClient{
		domains:  []types.DomainInfo{},
		statuses: map[string]*types.DomainStatus{},
	}

	fetcher := &OpenSearchFetcher{client: mock}
	resources, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, resources)
}

func TestMapOpenSearchDomain_DeletedFlag(t *testing.T) {
	status := types.DomainStatus{
		DomainName:    aws.String("deleted-domain"),
		ARN:           aws.String("arn:aws:es:us-east-1:123456789012:domain/deleted-domain"),
		EngineVersion: aws.String("OpenSearch_1.3"),
		Deleted:       aws.Bool(true),
		ClusterConfig: &types.ClusterConfig{
			InstanceType: types.OpenSearchPartitionInstanceTypeT3SmallSearch,
		},
	}

	result := mapOpenSearchDomain(status)

	assert.Equal(t, "deleted-domain", result.ID)
	assert.True(t, result.Deleted)
	assert.Equal(t, "OpenSearch_1.3", result.EngineVersion)
}
