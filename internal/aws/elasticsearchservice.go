package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticsearchservice"
)

type ESClient interface {
	ListDomainNames(ctx context.Context, params *elasticsearchservice.ListDomainNamesInput, optFns ...func(*elasticsearchservice.Options)) (*elasticsearchservice.ListDomainNamesOutput, error)
	DescribeElasticsearchDomains(ctx context.Context, params *elasticsearchservice.DescribeElasticsearchDomainsInput, optFns ...func(*elasticsearchservice.Options)) (*elasticsearchservice.DescribeElasticsearchDomainsOutput, error)
}

type ESResource struct {
	ID      string
	Name    string
	Version string
	Deleted bool
}

type ESFetcher struct {
	client ESClient
}

func NewESFetcher(client ESClient) *ESFetcher {
	return &ESFetcher{client: client}
}

func (f *ESFetcher) FetchAll(ctx context.Context) ([]ESResource, error) {
	listOut, err := f.client.ListDomainNames(ctx, &elasticsearchservice.ListDomainNamesInput{})
	if err != nil {
		return nil, err
	}

	var names []string
	for _, d := range listOut.DomainNames {
		if d.DomainName != nil {
			names = append(names, aws.ToString(d.DomainName))
		}
	}

	if len(names) == 0 {
		return []ESResource{}, nil
	}

	descOut, err := f.client.DescribeElasticsearchDomains(ctx, &elasticsearchservice.DescribeElasticsearchDomainsInput{
		DomainNames: names,
	})
	if err != nil {
		return nil, err
	}

	var results []ESResource
	for _, d := range descOut.DomainStatusList {
		results = append(results, mapESDomain(d))
	}
	return results, nil
}

func mapESDomain(d elasticsearchservice.ElasticsearchDomainStatus) ESResource {
	deleted := false
	if d.Deleted != nil {
		deleted = *d.Deleted
	}
	return ESResource{
		ID:      aws.ToString(d.ARN),
		Name:    aws.ToString(d.DomainName),
		Version: aws.ToString(d.ElasticsearchVersion),
		Deleted: deleted,
	}
}
