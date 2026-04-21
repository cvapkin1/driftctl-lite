package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type ECSCluster struct {
	ARN    string
	Name   string
	Status string
}

type ECSFetcher struct {
	client *ecs.Client
}

func NewECSFetcher(client *ecs.Client) *ECSFetcher {
	return &ECSFetcher{client: client}
}

func (f *ECSFetcher) FetchAll(ctx context.Context) ([]ECSCluster, error) {
	var clusters []ECSCluster

	listOut, err := f.client.ListClusters(ctx, &ecs.ListClustersInput{})
	if err != nil {
		return nil, err
	}
	if len(listOut.ClusterArns) == 0 {
		return clusters, nil
	}

	descOut, err := f.client.DescribeClusters(ctx, &ecs.DescribeClustersInput{
		Clusters: listOut.ClusterArns,
	})
	if err != nil {
		return nil, err
	}

	for _, c := range descOut.Clusters {
		clusters = append(clusters, mapECSCluster(c))
	}
	return clusters, nil
}

func mapECSCluster(c ecs.Cluster) ECSCluster {
	return ECSCluster{
		ARN:    aws.ToString(c.ClusterArn),
		Name:   aws.ToString(c.ClusterName),
		Status: aws.ToString(c.Status),
	}
}
