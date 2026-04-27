package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/budgets"
	"github.com/aws/aws-sdk-go-v2/service/budgets/types"
)

type BudgetResource struct {
	AccountID   string
	Name        string
	BudgetType  string
	LimitAmount string
	LimitUnit   string
	TimeUnit    string
}

type BudgetsClient interface {
	DescribeBudgets(ctx context.Context, params *budgets.DescribeBudgetsInput, optFns ...func(*budgets.Options)) (*budgets.DescribeBudgetsOutput, error)
}

type BudgetsFetcher struct {
	client    BudgetsClient
	accountID string
}

func NewBudgetsFetcher(client BudgetsClient, accountID string) *BudgetsFetcher {
	return &BudgetsFetcher{client: client, accountID: accountID}
}

func (f *BudgetsFetcher) FetchAll(ctx context.Context) ([]BudgetResource, error) {
	out, err := f.client.DescribeBudgets(ctx, &budgets.DescribeBudgetsInput{
		AccountId: aws.String(f.accountID),
	})
	if err != nil {
		return nil, err
	}

	resources := make([]BudgetResource, 0, len(out.Budgets))
	for _, b := range out.Budgets {
		resources = append(resources, mapBudget(b, f.accountID))
	}
	return resources, nil
}

func mapBudget(b types.Budget, accountID string) BudgetResource {
	r := BudgetResource{
		AccountID:  accountID,
		BudgetType: string(b.BudgetType),
		TimeUnit:   string(b.TimeUnit),
	}
	if b.BudgetName != nil {
		r.Name = *b.BudgetName
	}
	if b.BudgetLimit != nil {
		if b.BudgetLimit.Amount != nil {
			r.LimitAmount = *b.BudgetLimit.Amount
		}
		if b.BudgetLimit.Unit != nil {
			r.LimitUnit = *b.BudgetLimit.Unit
		}
	}
	return r
}
