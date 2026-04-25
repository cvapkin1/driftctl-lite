package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/redshift"
	"github.com/aws/aws-sdk-go-v2/service/redshift/types"
)

type RedshiftCluster struct {
	ID             string
	NodeType       string
	Status         string
	DBName         string
	NumberOfNodes  int32
}

type RedshiftClient interface {
	DescribeClusters(ctx context.Context, params *redshift.DescribeClustersInput, optFns ...func(*redshift.Options)) (*redshift.DescribeClustersOutput, error)
}

type RedshiftFetcher struct {
	client RedshiftClient
}

func NewRedshiftFetcher(client RedshiftClient) *RedshiftFetcher {
	return &RedshiftFetcher{client: client}
}

func (f *RedshiftFetcher) FetchAll(ctx context.Context) ([]RedshiftCluster, error) {
	var clusters []RedshiftCluster
	var marker *string

	for {
		out, err := f.client.DescribeClusters(ctx, &redshift.DescribeClustersInput{
			Marker: marker,
		})
		if err != nil {
			return nil, err
		}

		for _, c := range out.Clusters {
			clusters = append(clusters, mapRedshiftCluster(c))
		}

		if out.Marker == nil {
			break
		}
		marker = out.Marker
	}

	return clusters, nil
}

func mapRedshiftCluster(c types.Cluster) RedshiftCluster {
	return RedshiftCluster{
		ID:            aws.ToString(c.ClusterIdentifier),
		NodeType:      aws.ToString(c.NodeType),
		Status:        aws.ToString(c.ClusterStatus),
		DBName:        aws.ToString(c.DBName),
		NumberOfNodes: c.NumberOfNodes,
	}
}
