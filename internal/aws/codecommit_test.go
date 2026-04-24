package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"github.com/stretchr/testify/assert"
)

func TestMapRepository_Fields(t *testing.T) {
	out := &codecommit.GetRepositoryOutput{
		RepositoryMetadata: &types.RepositoryMetadata{
			RepositoryId:   aws.String("repo-id-123"),
			RepositoryName: aws.String("my-repo"),
			CloneUrlHttp:   aws.String("https://git-codecommit.us-east-1.amazonaws.com/v1/repos/my-repo"),
			Arn:            aws.String("arn:aws:codecommit:us-east-1:123456789012:my-repo"),
		},
	}

	repo := mapRepository(out)

	assert.Equal(t, "repo-id-123", repo.ID)
	assert.Equal(t, "my-repo", repo.Name)
	assert.Equal(t, "https://git-codecommit.us-east-1.amazonaws.com/v1/repos/my-repo", repo.CloneURLHTTP)
	assert.Equal(t, "arn:aws:codecommit:us-east-1:123456789012:my-repo", repo.ARN)
}

func TestMapRepository_EmptyFields(t *testing.T) {
	out := &codecommit.GetRepositoryOutput{
		RepositoryMetadata: &types.RepositoryMetadata{},
	}

	repo := mapRepository(out)

	assert.Equal(t, "", repo.ID)
	assert.Equal(t, "", repo.Name)
	assert.Equal(t, "", repo.CloneURLHTTP)
	assert.Equal(t, "", repo.ARN)
}
