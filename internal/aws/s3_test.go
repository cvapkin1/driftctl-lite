package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

typeS3Client struct {
	bttags    map[string][].Tag
}

 context.Context, _Input, _ ...func(*s3BucketsOutput, error) {
	return &s3.ListBucketsOutput{Buckets: m.buckets}, nil
}

func (m *mockS3Client) GetBucketLocation(ctx context.Context, in *s3.GetBucketLocationInput, _ ...func(*s3.Options)) (*s3.GetBucketLocationOutput, error) {
	return &s3.GetBucketLocationOutput{LocationConstraint: types.BucketLocationConstraintEuWest1}, nil
}

func (m *mockS3Client) GetBucketTagging(ctx context.Context, in *s3.GetBucketTaggingInput, _ ...func(*s3.Options)) (*s3.GetBucketTaggingOutput, error) {
	tags := m.tags[aws.ToString(in.Bucket)]
	return &s3.GetBucketTaggingOutput{TagSet: tags}, nil
}

func TestS3FetchAll_ReturnsBuckets(t *testing.T) {
	client := &mockS3Client{
		buckets: []types.Bucket{
			{Name: aws.String("my-bucket")},
		},
		tags: map[string][]types.Tag{
			"my-bucket": {{Key: aws.String("env"), Value: aws.String("prod")}},
		},
	}
	fetcher := NewS3Fetcher(client)
	buckets, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(buckets) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(buckets))
	}
	if buckets[0].Name != "my-bucket" {
		t.Errorf("expected name my-bucket, got %s", buckets[0].Name)
	}
	if buckets[0].Tags["env"] != "prod" {
		t.Errorf("expected tag env=prod")
	}
	if buckets[0].Region != "eu-west-1" {
		t.Errorf("expected region eu-west-1, got %s", buckets[0].Region)
	}
}
