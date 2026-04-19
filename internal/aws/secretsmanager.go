package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerClient interface {
	ListSecrets(ctx context.Context, params *secretsmanager.ListSecretsInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.ListSecretsOutput, error)
}

type SecretsManagerFetcher struct {
	client SecretsManagerClient
}

func NewSecretsManagerFetcher(client SecretsManagerClient) *SecretsManagerFetcher {
	return &SecretsManagerFetcher{client: client}
}

type SecretResource struct {
	ARN  string
	Name string
}

func (f *SecretsManagerFetcher) FetchAll(ctx context.Context) ([]SecretResource, error) {
	var secrets []SecretResource
	var nextToken *string

	for {
		resp, err := f.client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, s := range resp.SecretList {
			secrets = append(secrets, SecretResource{
				ARN:  aws.ToString(s.ARN),
				Name: aws.ToString(s.Name),
			})
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return secrets, nil
}
