package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/stretchr/testify/assert"
)

type mockAPIGatewayClient struct {
	output *apigateway.GetRestApisOutput
	err    error
}

func (m *mockAPIGatewayClient) GetRestApis(_ context.Context, _ *apigateway.GetRestApisInput, _ ...func(*apigateway.Options)) (*apigateway.GetRestApisOutput, error) {
	return m.output, m.err
}

func TestAPIGatewayFetchAll_ReturnsAPIs(t *testing.T) {
	mock := &mockAPIGatewayClient{
		output: &apigateway.GetRestApisOutput{
			Items: []apigateway.RestApi{
				{Id: aws.String("abc123"), Name: aws.String("my-api")},
				{Id: aws.String("def456"), Name: aws.String("other-api")},
			},
		},
	}
	fetcher := &APIGatewayFetcher{client: mock, region: "us-east-1"}
	result, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "abc123", result[0].ID)
	assert.Equal(t, "my-api", result[0].Name)
	assert.Equal(t, "arn:aws:apigateway:us-east-1::/restapis/abc123", result[0].ARN)
}

func TestAPIGatewayFetchAll_Empty(t *testing.T) {
	mock := &mockAPIGatewayClient{
		output: &apigateway.GetRestApisOutput{Items: []apigateway.RestApi{}},
	}
	fetcher := &APIGatewayFetcher{client: mock, region: "us-east-1"}
	result, err := fetcher.FetchAll(context.Background())
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestMapRestAPI_Fields(t *testing.T) {
	api := apigateway.RestApi{
		Id:   aws.String("xyz789"),
		Name: aws.String("test-api"),
	}
	res := mapRestAPI(api, "eu-west-1")
	assert.Equal(t, "xyz789", res.ID)
	assert.Equal(t, "test-api", res.Name)
	assert.Equal(t, "arn:aws:apigateway:eu-west-1::/restapis/xyz789", res.ARN)
}
