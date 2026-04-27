package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type CloudFrontResource struct {
	ID           string
	DomainName   string
	Status       string
	Enabled      bool
	PriceClass   string
	Comment      string
}

type CloudFrontAPI interface {
	ListDistributions(ctx context.Context, params *cloudfront.ListDistributionsInput, optFns ...func(*cloudfront.Options)) (*cloudfront.ListDistributionsOutput, error)
}

type CloudFrontFetcher struct {
	client CloudFrontAPI
}

func NewCloudFrontFetcher(client CloudFrontAPI) *CloudFrontFetcher {
	return &CloudFrontFetcher{client: client}
}

func (f *CloudFrontFetcher) FetchAll(ctx context.Context) ([]CloudFrontResource, error) {
	var resources []CloudFrontResource
	var marker *string

	for {
		out, err := f.client.ListDistributions(ctx, &cloudfront.ListDistributionsInput{
			Marker: marker,
		})
		if err != nil {
			return nil, err
		}
		if out.DistributionList == nil {
			break
		}
		for _, d := range out.DistributionList.Items {
			resources = append(resources, mapDistribution(d))
		}
		if !aws.ToBool(out.DistributionList.IsTruncated) {
			break
		}
		marker = out.DistributionList.NextMarker
	}
	return resources, nil
}

func mapDistribution(d types.DistributionSummary) CloudFrontResource {
	return CloudFrontResource{
		ID:         aws.ToString(d.Id),
		DomainName: aws.ToString(d.DomainName),
		Status:     aws.ToString(d.Status),
		Enabled:    aws.ToBool(d.Enabled),
		PriceClass: string(d.PriceClass),
		Comment:    aws.ToString(d.Comment),
	}
}
