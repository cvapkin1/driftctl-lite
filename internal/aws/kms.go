package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

type KMSKey struct {
	KeyID       string
	ARN         string
	State       string
	Description string
}

type KMSClient interface {
	ListKeys(ctx context.Context, params *kms.ListKeysInput, optFns ...func(*kms.Options)) (*kms.ListKeysOutput, error)
	DescribeKey(ctx context.Context, params *kms.DescribeKeyInput, optFns ...func(*kms.Options)) (*kms.DescribeKeyOutput, error)
}

type KMSFetcher struct {
	client KMSClient
}

func NewKMSFetcher(client KMSClient) *KMSFetcher {
	return &KMSFetcher{client: client}
}

func (f *KMSFetcher) FetchAll(ctx context.Context) ([]KMSKey, error) {
	var keys []KMSKey

	out, err := f.client.ListKeys(ctx, &kms.ListKeysInput{})
	if err != nil {
		return nil, err
	}

	for _, entry := range out.Keys {
		desc, err := f.client.DescribeKey(ctx, &kms.DescribeKeyInput{
			KeyId: entry.KeyId,
		})
		if err != nil {
			return nil, err
		}
		keys = append(keys, mapKMSKey(desc.KeyMetadata))
	}

	return keys, nil
}

func mapKMSKey(m *types.KeyMetadata) KMSKey {
	return KMSKey{
		KeyID:       aws.ToString(m.KeyId),
		ARN:         aws.ToString(m.Arn),
		State:       string(m.KeyState),
		Description: aws.ToString(m.Description),
	}
}
