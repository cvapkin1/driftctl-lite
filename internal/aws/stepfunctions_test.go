package aws_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"github.com/stretchr/testify/assert"
)

type mockSFNClient struct {
	machines []types.StateMachineListItem
}

func (m *mockSFNClient) ListStateMachines(ctx context.Context, input *sfn.ListStateMachinesInput, opts ...func(*sfn.Options)) (*sfn.ListStateMachinesOutput, error) {
	return &sfn.ListStateMachinesOutput{
		StateMachines: m.machines,
	}, nil
}

func TestStepFunctionsFetchAll_ReturnsMachines(t *testing.T) {
	machines := []types.StateMachineListItem{
		{Arn: aws.String("arn:aws:states:us-east-1:123456789012:stateMachine:MyMachine"), Name: aws.String("MyMachine")},
		{Arn: aws.String("arn:aws:states:us-east-1:123456789012:stateMachine:OtherMachine"), Name: aws.String("OtherMachine")},
	}

	result := mapStateMachinesForTest(machines)
	assert.Len(t, result, 2)
	assert.Equal(t, "arn:aws:states:us-east-1:123456789012:stateMachine:MyMachine", result[0].ARN)
	assert.Equal(t, "MyMachine", result[0].Name)
}

func TestStepFunctionsFetchAll_Empty(t *testing.T) {
	result := mapStateMachinesForTest([]types.StateMachineListItem{})
	assert.Empty(t, result)
}

func mapStateMachinesForTest(items []types.StateMachineListItem) []struct{ ARN, Name string } {
	var out []struct{ ARN, Name string }
	for _, item := range items {
		out = append(out, struct{ ARN, Name string }{
			ARN:  aws.ToString(item.Arn),
			Name: aws.ToString(item.Name),
		})
	}
	return out
}
