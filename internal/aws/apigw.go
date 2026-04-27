package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/aws/aws-sdk-go-v2/config"
)

type APIGatewayResource struct {
	ID   string
	Name string
	ARN  string
}

type apiGatewayClient interface {
	GetRestApis(ctx context.Context, params *apigateway.GetRestApisInput, optFns ...func(*apigateway.Options)) (*apigateway.GetRestApisOutput, error)
}

type APIGatewayFetcher struct {
	client apiGatewayClient
	region string
}

func NewAPIGatewayFetcher(ctx context.Context, region string) (*APIGatewayFetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &APIGatewayFetcher{
		client: apigateway.NewFromConfig(cfg),
		region: region,
	}, nil
}

func (f *APIGatewayFetcher) FetchAll(ctx context.Context) ([]APIGatewayResource, error) {
	var resources []APIGatewayResource
	var position *string

	for {
		out, err := f.client.GetRestApis(ctx, &apigateway.GetRestApisInput{
			Position: position,
		})
		if err != nil {
			return nil, err
		}
		for _, api := range out.Items {
			resources = append(resources, mapRestAPI(api, f.region))
		}
		if out.Position == nil {
			break
		}
		position = out.Position
	}
	return resources, nil
}

func mapRestAPI(api apigateway.RestApi, region string) APIGatewayResource {
	id := aws.ToString(api.Id)
	return APIGatewayResource{
		ID:   id,
		Name: aws.ToString(api.Name),
		ARN:  "arn:aws:apigateway:" + region + "::/restapis/" + id,
	}
}
