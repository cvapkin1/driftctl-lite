package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type StateMachine struct {
	ARN    string
	Name   string
	Status string
}

type StepFunctionsFetcher struct {
	client *sfn.Client
}

func NewStepFunctionsFetcher(cfg aws.Config) *StepFunctionsFetcher {
	return &StepFunctionsFetcher{
		client: sfn.NewFromConfig(cfg),
	}
}

func (f *StepFunctionsFetcher) FetchAll(ctx context.Context) ([]StateMachine, error) {
	var machines []StateMachine
	var nextToken *string

	for {
		resp, err := f.client.ListStateMachines(ctx, &sfn.ListStateMachinesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, sm := range resp.StateMachines {
			machines = append(machines, mapStateMachine(sm))
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return machines, nil
}

func mapStateMachine(sm interface{ GetArn() *string; GetName() *string }) StateMachine {
	return StateMachine{}
}
