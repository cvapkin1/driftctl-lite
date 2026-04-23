package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline"
	"github.com/aws/aws-sdk-go-v2/service/codepipeline/types"
)

type mockCodePipelineClient struct {
	pipelines []types.PipelineSummary
	err       error
}

func (m *mockCodePipelineClient) ListPipelines(ctx context.Context, params *codepipeline.ListPipelinesInput, optFns ...func(*codepipeline.Options)) (*codepipeline.ListPipelinesOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &codepipeline.ListPipelinesOutput{Pipelines: m.pipelines}, nil
}

func (m *mockCodePipelineClient) GetPipelineState(ctx context.Context, params *codepipeline.GetPipelineStateInput, optFns ...func(*codepipeline.Options)) (*codepipeline.GetPipelineStateOutput, error) {
	return &codepipeline.GetPipelineStateOutput{}, nil
}

func TestCodePipelineFetchAll_ReturnsPipelines(t *testing.T) {
	mock := &mockCodePipelineClient{
		pipelines: []types.PipelineSummary{
			{Name: aws.String("my-pipeline"), PipelineArn: aws.String("arn:aws:codepipeline:us-east-1:123456789012:my-pipeline")},
			{Name: aws.String("other-pipeline"), PipelineArn: aws.String("arn:aws:codepipeline:us-east-1:123456789012:other-pipeline")},
		},
	}
	fetcher := NewCodePipelineFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 pipelines, got %d", len(results))
	}
	if results[0].Name != "my-pipeline" {
		t.Errorf("expected name 'my-pipeline', got %s", results[0].Name)
	}
	if results[0].ARN != "arn:aws:codepipeline:us-east-1:123456789012:my-pipeline" {
		t.Errorf("unexpected ARN: %s", results[0].ARN)
	}
}

func TestCodePipelineFetchAll_Empty(t *testing.T) {
	mock := &mockCodePipelineClient{pipelines: []types.PipelineSummary{}}
	fetcher := NewCodePipelineFetcher(mock)
	results, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
