package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type mockLambdaClient struct {
	functions []types.FunctionConfiguration
}

func (m *mockLambdaClient) ListFunctions(_ context.Context, _ *lambda.ListFunctionsInput, _ ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error) {
	return &lambda.ListFunctionsOutput{Functions: m.functions}, nil
}

func TestLambdaFetchAll_ReturnsFunctions(t *testing.T) {
	mock := &mockLambdaClient{
		functions: []types.FunctionConfiguration{
			{FunctionName: aws.String("my-func"), Runtime: types.RuntimeNodejs18x, State: types.StateActive},
		},
	}
	fetcher := NewLambdaFetcherV2(mock)
	fns, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fns) != 1 {
		t.Fatalf("expected 1 function, got %d", len(fns))
	}
	if fns[0].Name != "my-func" {
		t.Errorf("expected name my-func, got %s", fns[0].Name)
	}
}

func TestLambdaFetchAll_Empty(t *testing.T) {
	mock := &mockLambdaClient{}
	fetcher := NewLambdaFetcherV2(mock)
	fns, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(fns) != 0 {
		t.Errorf("expected 0 functions, got %d", len(fns))
	}
}
