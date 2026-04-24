package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
)

type CodeCommitRepository struct {
	ID          string
	Name        string
	CloneURLHTTP string
	ARN         string
	State       string
}

type CodeCommitFetcher struct {
	client *codecommit.Client
}

func NewCodeCommitFetcher(client *codecommit.Client) *CodeCommitFetcher {
	return &CodeCommitFetcher{client: client}
}

func (f *CodeCommitFetcher) FetchAll(ctx context.Context) ([]CodeCommitRepository, error) {
	var repos []CodeCommitRepository
	var nextToken *string

	for {
		out, err := f.client.ListRepositories(ctx, &codecommit.ListRepositoriesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range out.Repositories {
			detail, err := f.client.GetRepository(ctx, &codecommit.GetRepositoryInput{
				RepositoryName: r.RepositoryName,
			})
			if err != nil {
				continue
			}
			repos = append(repos, mapRepository(detail.RepositoryMetadata))
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return repos, nil
}

func mapRepository(r *codecommit.GetRepositoryOutput) CodeCommitRepository {
	if r == nil {
		return CodeCommitRepository{}
	}
	return CodeCommitRepository{
		ID:           aws.ToString(r.RepositoryMetadata.RepositoryId),
		Name:         aws.ToString(r.RepositoryMetadata.RepositoryName),
		CloneURLHTTP: aws.ToString(r.RepositoryMetadata.CloneUrlHttp),
		ARN:          aws.ToString(r.RepositoryMetadata.Arn),
	}
}
