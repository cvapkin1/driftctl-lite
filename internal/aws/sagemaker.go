package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker"
	"github.com/aws/aws-sdk-go-v2/service/sagemaker/types"
)

type SageMakerModel struct {
	Name    string
	ARN     string
	Status  string
}

type sageMakerClient interface {
	ListModels(ctx context.Context, params *sagemaker.ListModelsInput, optFns ...func(*sagemaker.Options)) (*sagemaker.ListModelsOutput, error)
}

type SageMakerFetcher struct {
	client sageMakerClient
}

func NewSageMakerFetcher(client sageMakerClient) *SageMakerFetcher {
	return &SageMakerFetcher{client: client}
}

func (f *SageMakerFetcher) FetchAll(ctx context.Context) ([]SageMakerModel, error) {
	var models []SageMakerModel
	var nextToken *string

	for {
		out, err := f.client.ListModels(ctx, &sagemaker.ListModelsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, m := range out.Models {
			models = append(models, mapSageMakerModel(m))
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return models, nil
}

func mapSageMakerModel(m types.ModelSummary) SageMakerModel {
	return SageMakerModel{
		Name:   aws.ToString(m.ModelName),
		ARN:    aws.ToString(m.ModelArn),
		Status: "active",
	}
}
