package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/aws/aws-sdk-go-v2/service/kafka/types"
)

type MSKCluster struct {
	ARN          string
	Name         string
	State        string
	KafkaVersion string
	BrokerCount  int32
}

type MSKClient interface {
	ListClusters(ctx context.Context, params *kafka.ListClustersV2Input, optFns ...func(*kafka.Options)) (*kafka.ListClustersV2Output, error)
}

type MSKFetcher struct {
	client MSKClient
}

func NewMSKFetcher(client MSKClient) *MSKFetcher {
	return &MSKFetcher{client: client}
}

func (f *MSKFetcher) FetchAll(ctx context.Context) ([]MSKCluster, error) {
	var clusters []MSKCluster
	var nextToken *string

	for {
		out, err := f.client.ListClusters(ctx, &kafka.ListClustersV2Input{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, c := range out.ClusterInfoList {
			clusters = append(clusters, mapMSKCluster(c))
		}

		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}

	return clusters, nil
}

func mapMSKCluster(c types.Cluster) MSKCluster {
	cluster := MSKCluster{
		ARN:   aws.ToString(c.ClusterArn),
		Name:  aws.ToString(c.ClusterName),
		State: string(c.State),
	}
	if c.Provisioned != nil {
		cluster.KafkaVersion = aws.ToString(c.Provisioned.CurrentBrokerSoftwareInfo.KafkaVersion)
		cluster.BrokerCount = c.Provisioned.NumberOfBrokerNodes
	}
	return cluster
}
