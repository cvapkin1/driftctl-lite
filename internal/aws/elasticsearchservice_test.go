package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice/types"
	"github.com/stretchr/testify/assert"
)

type mockESClient struct {
	listOut *elasticsearchservice.ListDomainNamesOutput
	descOut *elasticsearchservice.DescribeElasticsearchDomainsOutput
}

func (m *mockESClient) ListDomainNames(ctx context.Context, params *elasticsearchservice.ListDomainNamesInput, optFns ...func(*elasticsearchservice.Options)) (*elasticsearchservice.ListDomainNamesOutput, error) {
	return m.listOut, nil
}

func (m *mockESClient) DescribeElasticsearchDomains(ctx context.Context, params *elasticsearchservice.DescribeElasticsearchDomainsInput, optFns ...func(*elasticsearchservice.Options)) (*elasticsearchservice.DescribeElasticsearchDomainsOutput, error) {
	return m.descOut, nil
}

func TestESFetchAll_ReturnsDomains(t *testing.T) {
	client := &mockESClient{
		listOut: &elasticsearchservice.ListDomainNamesOutput{
			DomainNames: []types.DomainInfo{{DomainName: aws.String("my-domain")}},
		},
		descOut: &elasticsearchservice.DescribeElasticsearchDomainsOutput{
			DomainStatusList: []elasticsearchservice.ElasticsearchDomainStatus{
				{
					ARN:                  aws.String("arn:aws:es:us-east-1:123456789012:domain/my-domain"),
					DomainName:           aws.String("my-domain"),
					ElasticsearchVersion: aws.String("7.10"),
					Deleted:              aws.Bool(false),
				},
			},
		},
	}
	fetcher := NewESFetcher(client)
	results, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "my-domain", results[0].Name)
	assert.Equal(t, "7.10", results[0].Version)
	assert.False(t, results[0].Deleted)
}

func TestESFetchAll_Empty(t *testing.T) {
	client := &mockESClient{
		listOut: &elasticsearchservice.ListDomainNamesOutput{
			DomainNames: []types.DomainInfo{},
		},
	}
	fetcher := NewESFetcher(client)
	results, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, results)
}
