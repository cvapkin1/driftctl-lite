package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/emr"
	"github.com/aws/aws-sdk-go-v2/service/emr/types"
)

type EMRCluster struct {
	ID     string
	Name   string
	State  string
	ARN    string
}

type emrClient interface {
	ListClusters(ctx context.Context, params *emr.ListClustersInput, optFns ...func(*emr.Options)) (*emr.ListClustersOutput, error)
	DescribeCluster(ctx context.Context, params *emr.DescribeClusterInput, optFns ...func(*emr.Options)) (*emr.DescribeClusterOutput, error)
}

type EMRFetcher struct {
	client emrClient
}

func NewEMRFetcher(client emrClient) *EMRFetcher {
	return &EMRFetcher{client: client}
}

func (f *EMRFetcher) FetchAll(ctx context.Context) ([]EMRCluster, error) {
	var clusters []EMRCluster
	var marker *string

	for {
		out, err := f.client.ListClusters(ctx, &emr.ListClustersInput{
			Marker: marker,
		})
		if err != nil {
			return nil, err
		}

		for _, summary := range out.Clusters {
			desc, err := f.client.DescribeCluster(ctx, &emr.DescribeClusterInput{
				ClusterId: summary.Id,
			})
			if err != nil {
				return nil, err
			}
			clusters = append(clusters, mapEMRCluster(desc.Cluster))
		}

		if out.Marker == nil {
			break
		}
		marker = out.Marker
	}

	return clusters, nil
}

func mapEMRCluster(c *types.Cluster) EMRCluster {
	if c == nil {
		return EMRCluster{}
	}
	return EMRCluster{
		ID:    aws.ToString(c.Id),
		Name:  aws.ToString(c.Name),
		State: string(c.Status.State),
		ARN:   aws.ToString(c.ClusterArn),
	}
}
