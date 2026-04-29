package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type ElasticIPResource struct {
	AllocationID string
	PublicIP     string
	InstanceID   string
	Domain       string
	Associated   bool
}

type elasticIPClient interface {
	DescribeAddresses(ctx context.Context, params *ec2.DescribeAddressesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeAddressesOutput, error)
}

type ElasticIPFetcher struct {
	client elasticIPClient
}

func NewElasticIPFetcher(client elasticIPClient) *ElasticIPFetcher {
	return &ElasticIPFetcher{client: client}
}

func (f *ElasticIPFetcher) FetchAll(ctx context.Context) ([]ElasticIPResource, error) {
	out, err := f.client.DescribeAddresses(ctx, &ec2.DescribeAddressesInput{})
	if err != nil {
		return nil, err
	}

	results := make([]ElasticIPResource, 0, len(out.Addresses))
	for _, addr := range out.Addresses {
		results = append(results, mapElasticIP(addr))
	}
	return results, nil
}

func mapElasticIP(addr types.Address) ElasticIPResource {
	res := ElasticIPResource{
		AllocationID: aws.ToString(addr.AllocationId),
		PublicIP:     aws.ToString(addr.PublicIp),
		InstanceID:   aws.ToString(addr.InstanceId),
		Domain:       string(addr.Domain),
		Associated:   addr.AssociationId != nil,
	}
	return res
}
