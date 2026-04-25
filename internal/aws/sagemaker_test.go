package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
	"github.com/stretchr/testify/assert"
)

type mockSageMakerClient struct {
	output *sagemaker.ListModelsOutput
	err    error
}

func (m *mockSageMakerClient) ListModels(ctx context.Context, params *sagemaker.ListModelsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListModelsOutput, error) {
	return m.output, m.err
}

func TestSageMakerFetchAll_ReturnsModels(t *testing.T) {
	mock := &mockSageMakerClient{
		output: &sagemaker.ListModelsOutput{
			Models: []types.ModelSummary{
				{ModelName: aws.String("model-a"), ModelArn: aws.String("arn:aws:sagemaker:us-east-1:123:model/model-a")},
				{ModelName: aws.String("model-b"), ModelArn: aws.String("arn:aws:sagemaker:us-east-1:123:model/model-b")},
			},
		},
	}
	fetcher := NewSageMakerFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "model-a", results[0].Name)
	assert.Equal(t, "arn:aws:sagemaker:us-east-1:123:model/model-a", results[0].ARN)
}

func TestSageMakerFetchAll_Empty(t *testing.T) {
	mock := &mockSageMakerClient{
		output: &sagemaker.ListModelsOutput{},
	}
	fetcher := NewSageMakerFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, results)
}
