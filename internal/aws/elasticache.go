package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/elasticache"
	"github.com/aws/aws-sdk-go-v2/config"
)

type ElastiCacheCluster struct {
	ClusterID string
	Engine    string
	Status    string
	NodeType  string
}

type ElastiCacheFetcher struct {
	client *elasticache.Client
}

func NewElastiCacheFetcher(regionElastiCacheFetcher, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	retetcher{client: elasticache.NewFromConfig(cfg)}, nil
}

func (f *ElastiCacheFetcher) FetchAll(ctx context.Context) ([]ElastiCacheCluster, error) {
	var clusters []ElastiCacheCluster
	paginator := elasticache.NewDescribeCacheClustersPaginator(f.client, &elasticache.DescribeCacheClustersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, c := range page.CacheClusters {
			clusters = append(clusters, mapElastiCacheCluster(c))
		}
	}
	return clusters, nil
}

func mapElastiCacheCluster(c elasticachetypes.CacheCluster) ElastiCacheCluster {
	return ElastiCacheCluster{
		ClusterID: aws.ToString(c.CacheClusterId),
		Engine:    aws.ToString(c.Engine),
		Status:    aws.ToString(c.CacheClusterStatus),
		NodeType:  aws.ToString(c.CacheNodeType),
	}
}
