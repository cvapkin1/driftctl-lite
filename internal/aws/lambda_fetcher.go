package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type lambdaClient interface {
	ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}

type LambdaFetcherV2 struct {
	client lambdaClient
}

func NewLambdaFetcherV2(client lambdaClient) *LambdaFetcherV2 {
	return &LambdaFetcherV2{client: client}
}

func (f *LambdaFetcherV2) FetchAll(ctx context.Context) ([]LambdaFunction, error) {
	var functions []LambdaFunction
	var marker *string
	for {
		out, err := f.client.ListFunctions(ctx, &lambda.ListFunctionsInput{Marker: marker})
		if err != nil {
			return nil, err
		}
		for _, fn := range out.Functions {
			functions = append(functions, mapLambdaConfig(fn))
		}
		if out.NextMarker == nil {
			break
		}
		marker = out.NextMarker
	}
	return functions, nil
}

func mapLambdaConfig(fn types.FunctionConfiguration) LambdaFunction {
	state := "active"
	if fn.State != "" {
		state = string(fn.State)
	}
	name := ""
	if fn.FunctionName != nil {
		name = *fn.FunctionName
	}
	runtime := string(fn.Runtime)
	return LambdaFunction{
		Name:    name,
		Runtime: runtime,
		State:   state,
	}
}
