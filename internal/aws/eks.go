package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

type EKSCluster struct {
	Name    string
	ARN     string
	Status  string
	Version string
}

type EKSFetcherAPI interface {
	ListClusters(ctx context.Context, params *eks.ListClustersInput, optFns ...func(*eks.Options)) (*eks.ListClustersOutput, error)
	DescribeCluster(ctx context.Context, params *eks.DescribeClusterInput, optFns ...func(*eks.Options)) (*eks.DescribeClusterOutput, error)
}

type EKSFetcher struct {
	client EKSFetcherAPI
}

func NewEKSFetcher(client EKSFetcherAPI) *EKSFetcher {
	return &EKSFetcher{client: client}
}

func (f *EKSFetcher) FetchAll(ctx context.Context) ([]EKSCluster, error) {
	listOut, err := f.client.ListClusters(ctx, &eks.ListClustersInput{})
	if err != nil {
		return nil, err
	}

	var clusters []EKSCluster
	for _, name := range listOut.Clusters {
		descOut, err := f.client.DescribeCluster(ctx, &eks.DescribeClusterInput{
			Name: aws.String(name),
		})
		if err != nil {
			return nil, err
		}
		if descOut.Cluster != nil {
			clusters = append(clusters, mapCluster(descOut.Cluster))
		}
	}
	return clusters, nil
}

func mapCluster(c *eksTypes.Cluster) EKSCluster {
	return EKSCluster{
		Name:    aws.ToString(c.Name),
		ARN:     aws.ToString(c.Arn),
		Status:  string(c.Status),
		Version: aws.ToString(c.Version),
	}
}
