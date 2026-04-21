package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/opensearch"
)

type OpenSearchDomain struct {
	DomainName    string
	ARN           string
	EngineVersion string
	Endpoint      string
	Deleted       bool
}

type openSearchClient interface {
	ListDomainNames(ctx context.Context, params *opensearch.ListDomainNamesInput, optFns ...func(*opensearch.Options)) (*opensearch.ListDomainNamesOutput, error)
	DescribeDomains(ctx context.Context, params *opensearch.DescribeDomainsInput, optFns ...func(*opensearch.Options)) (*opensearch.DescribeDomainsOutput, error)
}

type OpenSearchFetcher struct {
	client openSearchClient
}

func NewOpenSearchFetcher(client openSearchClient) *OpenSearchFetcher {
	return &OpenSearchFetcher{client: client}
}

func (f *OpenSearchFetcher) FetchAll(ctx context.Context) ([]OpenSearchDomain, error) {
	listOut, err := f.client.ListDomainNames(ctx, &opensearch.ListDomainNamesInput{})
	if err != nil {
		return nil, err
	}

	if len(listOut.DomainNames) == 0 {
		return []OpenSearchDomain{}, nil
	}

	names := make([]string, 0, len(listOut.DomainNames))
	for _, d := range listOut.DomainNames {
		if d.DomainName != nil {
			names = append(names, aws.ToString(d.DomainName))
		}
	}

	descOut, err := f.client.DescribeDomains(ctx, &opensearch.DescribeDomainsInput{
		DomainNames: names,
	})
	if err != nil {
		return nil, err
	}

	domains := make([]OpenSearchDomain, 0, len(descOut.DomainStatusList))
	for _, s := range descOut.DomainStatusList {
		domains = append(domains, mapOpenSearchDomain(s))
	}
	return domains, nil
}

func mapOpenSearchDomain(s opensearch.DomainStatus) OpenSearchDomain {
	endpoint := ""
	if s.Endpoint != nil {
		endpoint = aws.ToString(s.Endpoint)
	}
	deleted := false
	if s.Deleted != nil {
		deleted = aws.ToBool(s.Deleted)
	}
	return OpenSearchDomain{
		DomainName:    aws.ToString(s.DomainName),
		ARN:           aws.ToString(s.ARN),
		EngineVersion: aws.ToString(s.EngineVersion),
		Endpoint:      endpoint,
		Deleted:       deleted,
	}
}
