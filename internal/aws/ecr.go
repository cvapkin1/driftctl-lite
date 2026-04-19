package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

type ECRRepository struct {
	Name string
	ARN  string
	URI  string
}

type ECRFetcherAPI interface {
	DescribeRepositories(ctx context.Context, params *ecr.DescribeRepositoriesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error)
}

type ECRFetcher struct {
	client ECRFetcherAPI
}

func NewECRFetcher(client ECRFetcherAPI) *ECRFetcher {
	return &ECRFetcher{client: client}
}

func (f *ECRFetcher) FetchAll(ctx context.Context) ([]ECRRepository, error) {
	var repos []ECRRepository
	var nextToken *string

	for {
		out, err := f.client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range out.Repositories {
			repos = append(repos, ECRRepository{
				Name: aws.ToString(r.RepositoryName),
				ARN:  aws.ToString(r.RepositoryArn),
				URI:  aws.ToString(r.RepositoryUri),
			})
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return repos, nil
}
