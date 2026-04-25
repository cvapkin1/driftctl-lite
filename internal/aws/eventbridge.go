package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

type EventBridgeRule struct {
	Name        string
	ARN         string
	State       string
	Description string
	EventBusName string
}

type EventBridgeClient interface {
	ListRules(ctx context.Context, params *eventbridge.ListRulesInput, optFns ...func(*eventbridge.Options)) (*eventbridge.ListRulesOutput, error)
}

type EventBridgeFetcher struct {
	client EventBridgeClient
}

func NewEventBridgeFetcher(client EventBridgeClient) *EventBridgeFetcher {
	return &EventBridgeFetcher{client: client}
}

func (f *EventBridgeFetcher) FetchAll(ctx context.Context) ([]EventBridgeRule, error) {
	var rules []EventBridgeRule
	var nextToken *string

	for {
		out, err := f.client.ListRules(ctx, &eventbridge.ListRulesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, r := range out.Rules {
			rules = append(rules, mapEventBridgeRule(r))
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return rules, nil
}

func mapEventBridgeRule(r types.Rule) EventBridgeRule {
	return EventBridgeRule{
		Name:         aws.ToString(r.Name),
		ARN:          aws.ToString(r.Arn),
		State:        string(r.State),
		Description:  aws.ToString(r.Description),
		EventBusName: aws.ToString(r.EventBusName),
	}
}
