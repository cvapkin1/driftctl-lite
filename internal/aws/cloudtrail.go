package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
)

type CloudTrailClient interface {
	DescribeTrails(ctx context.Context, params *cloudtrail.DescribeTrailsInput, optFns ...func(*cloudtrail.Options)) (*cloudtrail.DescribeTrailsOutput, error)
}

type CloudTrailFetcher struct {
	client CloudTrailClient
}

type CloudTrailResource struct {
	Name      string
	ARN       string
	S3Bucket  string
	IsLogging bool
}

func NewCloudTrailFetcher(client CloudTrailClient) *CloudTrailFetcher {
	return &CloudTrailFetcher{client: client}
}

func (f *CloudTrailFetcher) FetchAll(ctx context.Context) ([]CloudTrailResource, error) {
	out, err := f.client.DescribeTrails(ctx, &cloudtrail.DescribeTrailsInput{
		IncludeShadowTrails: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}

	var resources []CloudTrailResource
	for _, trail := range out.TrailList {
		resources = append(resources, mapTrail(trail))
	}
	return resources, nil
}

func mapTrail(t types.Trail) CloudTrailResource {
	res := CloudTrailResource{}
	if t.Name != nil {
		res.Name = *t.Name
	}
	if t.TrailARN != nil {
		res.ARN = *t.TrailARN
	}
	if t.S3BucketName != nil {
		res.S3Bucket = *t.S3BucketName
	}
	return res
}
