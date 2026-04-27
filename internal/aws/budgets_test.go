package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/budgets"
	"github.com/aws/aws-sdk-go-v2/service/budgets/types"
	"github.com/stretchr/testify/assert"
)

type mockBudgetsClient struct {
	output *budgets.DescribeBudgetsOutput
	err    error
}

func (m *mockBudgetsClient) DescribeBudgets(ctx context.Context, params *budgets.DescribeBudgetsInput, optFns ...func(*budgets.Options)) (*budgets.DescribeBudgetsOutput, error) {
	return m.output, m.err
}

func TestBudgetsFetchAll_ReturnsBudgets(t *testing.T) {
	mock := &mockBudgetsClient{
		output: &budgets.DescribeBudgetsOutput{
			Budgets: []types.Budget{
				{
					BudgetName: aws.String("monthly-cost"),
					BudgetType: types.BudgetTypeCost,
					TimeUnit:   types.TimeUnitMonthly,
					BudgetLimit: &types.Spend{
						Amount: aws.String("1000"),
						Unit:   aws.String("USD"),
					},
				},
			},
		},
	}

	fetcher := NewBudgetsFetcher(mock, "123456789012")
	results, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "monthly-cost", results[0].Name)
	assert.Equal(t, "COST", results[0].BudgetType)
	assert.Equal(t, "1000", results[0].LimitAmount)
	assert.Equal(t, "USD", results[0].LimitUnit)
	assert.Equal(t, "MONTHLY", results[0].TimeUnit)
	assert.Equal(t, "123456789012", results[0].AccountID)
}

func TestBudgetsFetchAll_Empty(t *testing.T) {
	mock := &mockBudgetsClient{
		output: &budgets.DescribeBudgetsOutput{
			Budgets: []types.Budget{},
		},
	}

	fetcher := NewBudgetsFetcher(mock, "123456789012")
	results, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, results)
}
