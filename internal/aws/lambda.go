package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LambdaFunction struct {
	Name    string
	Runtime string
	State   string
}

type LambdaFetcher struct {
	client *lambda.Client
}

func NewLambdaFetcher(cfg aws.Config) *LambdaFetcher {
	return &LambdaFetcher{client: lambda.NewFromConfig(cfg)}
}

func (f *LambdaFetcher) FetchAll(ctx context.Context) ([]LambdaFunction, error) {
	var functions []LambdaFunction
	paginator := lambda.NewListFunctionsPaginator(f.client, &lambda.ListFunctionsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, fn := range page.Functions {
			functions = append(functions, mapLambda(fn))
		}
	}
	return functions, nil
}

func mapLambda(fn types.FunctionConfiguration) LambdaFunction {
	name := ""
	if fn.FunctionName != nil {
		name = *fn.FunctionName
	}
	return LambdaFunction{
		Name:    name,
		Runtime: string(fn.Runtime),
		State:   string(fn.State),
	}
}
