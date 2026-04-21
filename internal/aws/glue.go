package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/glue"
	"github.com/aws/aws-sdk-go-v2/service/glue/types"
)

type GlueJob struct {
	Name        string
	Role        string
	GlueVersion string
	State       string
}

type GlueFetcherAPI interface {
	GetJobs(ctx context.Context, params *glue.GetJobsInput, optFns ...func(*glue.Options)) (*glue.GetJobsOutput, error)
}

type GlueFetcher struct {
	client GlueFetcherAPI
}

func NewGlueFetcher(client GlueFetcherAPI) *GlueFetcher {
	return &GlueFetcher{client: client}
}

func (f *GlueFetcher) FetchAll(ctx context.Context) ([]GlueJob, error) {
	var jobs []GlueJob
	var nextToken *string

	for {
		out, err := f.client.GetJobs(ctx, &glue.GetJobsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, j := range out.Jobs {
			jobs = append(jobs, mapGlueJob(j))
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return jobs, nil
}

func mapGlueJob(j types.Job) GlueJob {
	name := aws.ToString(j.Name)
	role := aws.ToString(j.Role)
	glueVersion := aws.ToString(j.GlueVersion)
	return GlueJob{
		Name:        name,
		Role:        role,
		GlueVersion: glueVersion,
		State:       "ACTIVE",
	}
}
