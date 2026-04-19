package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBTable struct {
	Name   string
	ARN    string
	Status string
}

type DynamoDBFetcher struct {
	client *dynamodb.Client
}

func NewDynamoDBFetcher(client *dynamodb.Client) *DynamoDBFetcher {
	return &DynamoDBFetcher{client: client}
}

func (f *DynamoDBFetcher) FetchAll(ctx context.Context) ([]DynamoDBTable, error) {
	var tables []DynamoDBTable
	var lastEvaluated *string

	for {
		resp, err := f.client.ListTables(ctx, &dynamodb.ListTablesInput{
			ExclusiveStartTableName: lastEvaluated,
		})
		if err != nil {
			return nil, err
		}

		for _, name := range resp.TableNames {
			desc, err := f.client.DescribeTable(ctx, &dynamodb.DescribeTableInput{
				TableName: aws.String(name),
			})
			if err != nil {
				return nil, err
			}
			tables = append(tables, DynamoDBTable{
				Name:   name,
				ARN:    aws.ToString(desc.Table.TableArn),
				Status: string(desc.Table.TableStatus),
			})
		}

		if resp.LastEvaluatedTableName == nil {
			break
		}
		lastEvaluated = resp.LastEvaluatedTableName
	}

	return tables, nil
}
