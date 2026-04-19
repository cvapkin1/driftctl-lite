package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Bucket struct {
	ID     string
	Name   string
	Region string
	Tags   map[string]string
}

type S3API interface {
	ListBuckets(ctx context.Context, params *s3.ListBucketsInput, optFns ...func(*s3.Options)) (*s3.ListBucketsOutput, error)
	GetBucketLocation(ctx context.Context, params *s3.GetBucketLocationInput, optFns ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error)
	GetBucketTagging(ctx context.Context, params *s3.GetBucketTaggingInput, optFns ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error)
}

type S3Fetcher struct {
	client S3API
}

func NewS3Fetcher(client S3API) *S3Fetcher {
	return &S3Fetcher{client: client}
}

func (f *S3Fetcher) FetchAll(ctx context.Context) ([]S3Bucket, error) {
	out, err := f.client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var buckets []S3Bucket
	for _, b := range out.Buckets {
		name := aws.ToString(b.Name)

		loc, err := f.client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{Bucket: &name})
		region := "us-east-1"
		if err == nil && string(loc.LocationConstraint) != "" {
			region = string(loc.LocationConstraint)
		}

		tags := map[string]string{}
		tagOut, err := f.client.GetBucketTagging(ctx, &s3.GetBucketTaggingInput{Bucket: &name})
		if err == nil {
			for _, t := range tagOut.TagSet {
				tags[aws.ToString(t.Key)] = aws.ToString(t.Value)
			}
		}

		buckets = append(buckets, S3Bucket{
			ID:     name,
			Name:   name,
			Region: region,
			Tags:   tags,
		})
	}
	return buckets, nil
}
