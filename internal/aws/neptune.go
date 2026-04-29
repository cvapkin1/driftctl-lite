package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/neptune"
	"github.com/aws/aws-sdk-go-v2/service/neptune/types"
)

type NeptuneCluster struct {
	ClusterID     string
	ARN           string
	Status        string
	EngineVersion string
}

type neptuneClient interface {
	DescribeDBClusters(ctx context.Context, params *neptune.DescribeDBClustersInput, optFns ...func(*neptune.Options)) (*neptune.DescribeDBClustersOutput, error)
}

type NeptuneFetcher struct {
	client neptuneClient
}

func NewNeptuneFetcher(client neptuneClient) *NeptuneFetcher {
	return &NeptuneFetcher{client: client}
}

func (f *NeptuneFetcher) FetchAll(ctx context.Context) ([]NeptuneCluster, error) {
	var clusters []NeptuneCluster
	var marker *string

	for {
		out, err := f.client.DescribeDBClusters(ctx, &neptune.DescribeDBClustersInput{
			Filters: []types.Filter{
				{
					Name:   aws.String("engine"),
					Values: []string{"neptune"},
				},
			},
			Marker: marker,
		})
		if err != nil {
			return nil, err
		}

		for _, c := range out.DBClusters {
			clusters = append(clusters, mapNeptuneCluster(c))
		}

		if out.Marker == nil {
			break
		}
		marker = out.Marker
	}

	return clusters, nil
}

func mapNeptuneCluster(c types.DBCluster) NeptuneCluster {
	return NeptuneCluster{
		ClusterID:     aws.ToString(c.DBClusterIdentifier),
		ARN:           aws.ToString(c.DBClusterArn),
		Status:        aws.ToString(c.Status),
		EngineVersion: aws.ToString(c.EngineVersion),
	}
}
