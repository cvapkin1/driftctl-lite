package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/athena"
	"github.com/aws/aws-sdk-go-v2/service/athena/types"
)

type AthenaWorkgroup struct {
	Name   string
	State  string
	Engine string
}

type athenaClient interface {
	ListWorkGroups(ctx context.Context, params *athena.ListWorkGroupsInput, optFns ...func(*athena.Options)) (*athena.ListWorkGroupsOutput, error)
	GetWorkGroup(ctx context.Context, params *athena.GetWorkGroupInput, optFns ...func(*athena.Options)) (*athena.GetWorkGroupOutput, error)
}

type AthenaFetcher struct {
	client athenaClient
}

func NewAthenaFetcher(client athenaClient) *AthenaFetcher {
	return &AthenaFetcher{client: client}
}

func (f *AthenaFetcher) FetchAll(ctx context.Context) ([]AthenaWorkgroup, error) {
	var workgroups []AthenaWorkgroup

	listOut, err := f.client.ListWorkGroups(ctx, &athena.ListWorkGroupsInput{})
	if err != nil {
		return nil, err
	}

	for _, summary := range listOut.WorkGroups {
		getOut, err := f.client.GetWorkGroup(ctx, &athena.GetWorkGroupInput{
			WorkGroup: summary.Name,
		})
		if err != nil {
			continue
		}
		workgroups = append(workgroups, mapWorkgroup(getOut.WorkGroup))
	}

	return workgroups, nil
}

func mapWorkgroup(wg *types.WorkGroup) AthenaWorkgroup {
	engineVersion := ""
	if wg.Configuration != nil &&
		wg.Configuration.EngineVersion != nil &&
		wg.Configuration.EngineVersion.SelectedEngineVersion != nil {
		engineVersion = aws.ToString(wg.Configuration.EngineVersion.SelectedEngineVersion)
	}
	return AthenaWorkgroup{
		Name:   aws.ToString(wg.Name),
		State:  string(wg.State),
		Engine: engineVersion,
	}
}
