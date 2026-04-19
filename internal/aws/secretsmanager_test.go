package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

type mockSMClient struct {
	output *secretsmanager.ListSecretsOutput
	err    error
}

func (m *mockSMClient) ListSecrets(ctx context.Context, params *secretsmanager.ListSecretsInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.ListSecretsOutput, error) {
	return m.output, m.err
}

func TestSecretsManagerFetchAll_ReturnsSecrets(t *testing.T) {
	mock := &mockSMClient{
		output: &secretsmanager.ListSecretsOutput{
			SecretList: []types.SecretListEntry{
				{ARN: aws.String("arn:aws:secretsmanager:us-east-1:123456789012:secret:my-secret"), Name: aws.String("my-secret")},
				{ARN: aws.String("arn:aws:secretsmanager:us-east-1:123456789012:secret:other-secret"), Name: aws.String("other-secret")},
			},
		},
	}

	fetcher := NewSecretsManagerFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 secrets, got %d", len(results))
	}
	if results[0].Name != "my-secret" {
		t.Errorf("expected my-secret, got %s", results[0].Name)
	}
}

func TestSecretsManagerFetchAll_Empty(t *testing.T) {
	mock := &mockSMClient{
		output: &secretsmanager.ListSecretsOutput{},
	}

	fetcher := NewSecretsManagerFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 secrets, got %d", len(results))
	}
}
