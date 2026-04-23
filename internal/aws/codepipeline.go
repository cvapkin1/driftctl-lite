package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

type CodePipelineResource struct {
	Name    string
	ARN     string
	Status  string
}

type codePipelineClient interface {
	ListPipelines(ctx context.Context, params *codepipeline.ListPipelinesInput, optFns ...func(*codepipeline.Options)) (*codepipeline.ListPipelinesOutput, error)
	GetPipelineState(ctx context.Context, params *codepipeline.GetPipelineStateInput, optFns ...func(*codepipeline.Options)) (*codepipeline.GetPipelineStateOutput, error)
}

type CodePipelineFetcher struct {
	client codePipelineClient
}

func NewCodePipelineFetcher(client codePipelineClient) *CodePipelineFetcher {
	return &CodePipelineFetcher{client: client}
}

func (f *CodePipelineFetcher) FetchAll(ctx context.Context) ([]CodePipelineResource, error) {
	var resources []CodePipelineResource
	var nextToken *string

	for {
		out, err := f.client.ListPipelines(ctx, &codepipeline.ListPipelinesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, p := range out.Pipelines {
			resources = append(resources, mapPipeline(p))
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return resources, nil
}

func mapPipeline(p types.PipelineSummary) CodePipelineResource {
	return CodePipelineResource{
		Name:   aws.ToString(p.Name),
		ARN:    aws.ToString(p.PipelineArn),
		Status: "ACTIVE",
	}
}
