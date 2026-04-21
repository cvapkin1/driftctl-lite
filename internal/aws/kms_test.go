package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockKMSClient struct {
	listOut    *kms.ListKeysOutput
	describeOut *kms.DescribeKeyOutput
	err        error
}

func (m *mockKMSClient) ListKeys(ctx context.Context, params *kms.ListKeysInput, optFns ...func(*kms.Options)) (*kms.ListKeysOutput, error) {
	return m.listOut, m.err
}

func (m *mockKMSClient) DescribeKey(ctx context.Context, params *kms.DescribeKeyInput, optFns ...func(*kms.Options)) (*kms.DescribeKeyOutput, error) {
	return m.describeOut, m.err
}

func TestKMSFetchAll_ReturnsKeys(t *testing.T) {
	keyID := "arn:aws:kms:us-east-1:123456789012:key/abc-123"
	mock := &mockKMSClient{
		listOut: &kms.ListKeysOutput{
			Keys: []types.KeyListEntry{
				{KeyId: aws.String("abc-123")},
			},
		},
		describeOut: &kms.DescribeKeyOutput{
			KeyMetadata: &types.KeyMetadata{
				KeyId:       aws.String("abc-123"),
				Arn:         aws.String(keyID),
				KeyState:    types.KeyStateEnabled,
				Description: aws.String("my key"),
			},
		},
	}

	fetcher := NewKMSFetcher(mock)
	keys, err := fetcher.FetchAll(context.Background())

	require.NoError(t, err)
	assert.Len(t, keys, 1)
	assert.Equal(t, "abc-123", keys[0].KeyID)
	assert.Equal(t, keyID, keys[0].ARN)
	assert.Equal(t, "Enabled", keys[0].State)
	assert.Equal(t, "my key", keys[0].Description)
}

func TestKMSFetchAll_Empty(t *testing.T) {
	mock := &mockKMSClient{
		listOut: &kms.ListKeysOutput{Keys: []types.KeyListEntry{}},
	}

	fetcher := NewKMSFetcher(mock)
	keys, err := fetcher.FetchAll(context.Background())

	require.NoError(t, err)
	assert.Empty(t, keys)
}
