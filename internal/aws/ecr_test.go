package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type mockECRClient struct {
	output *ecr.DescribeRepositoriesOutput
	err    error
}

func (m *mockECRClient) DescribeRepositories(ctx context.Context, params *ecr.DescribeRepositoriesInput, optFns ...func(*ecr.Options)) (*ecr.DescribeRepositoriesOutput, error) {
	return m.output, m.err
}

func TestECRFetchAll_ReturnsRepositories(t *testing.T) {
	mock := &mockECRClient{
		output: &ecr.DescribeRepositoriesOutput{
			Repositories: []types.Repository{
				{
					RepositoryName: aws.String("my-repo"),
					RepositoryArn:  aws.String("arn:aws:ecr:us-east-1:123456789012:repository/my-repo"),
					RepositoryUri:  aws.String("123456789012.dkr.ecr.us-east-1.amazonaws.com/my-repo"),
				},
			},
		},
	}

	fetcher := NewECRFetcher(mock)
	repos, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected 1 repo, got %d", len(repos))
	}
	if repos[0].Name != "my-repo" {
		t.Errorf("expected name 'my-repo', got %s", repos[0].Name)
	}
}

func TestECRFetchAll_Empty(t *testing.T) {
	mock := &mockECRClient{
		output: &ecr.DescribeRepositoriesOutput{},
	}

	fetcher := NewECRFetcher(mock)
	repos, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(repos) != 0 {
		t.Errorf("expected 0 repos, got %d", len(repos))
	}
}
