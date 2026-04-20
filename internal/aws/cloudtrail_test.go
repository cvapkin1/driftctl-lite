package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCloudTrailClient struct {
	mock.Mock
}

func (m *mockCloudTrailClient) DescribeTrails(ctx context.Context, params *cloudtrail.DescribeTrailsInput, optFns ...func(*cloudtrail.Options)) (*cloudtrail.DescribeTrailsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*cloudtrail.DescribeTrailsOutput), args.Error(1)
}

func TestCloudTrailFetchAll_ReturnsTrails(t *testing.T) {
	mockClient := new(mockCloudTrailClient)
	mockClient.On("DescribeTrails", mock.Anything, mock.Anything).Return(&cloudtrail.DescribeTrailsOutput{
		TrailList: []types.Trail{
			{
				Name:         aws.String("my-trail"),
				TrailARN:     aws.String("arn:aws:cloudtrail:us-east-1:123456789012:trail/my-trail"),
				S3BucketName: aws.String("my-trail-bucket"),
			},
		},
	}, nil)

	fetcher := NewCloudTrailFetcher(mockClient)
	results, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "my-trail", results[0].Name)
	assert.Equal(t, "arn:aws:cloudtrail:us-east-1:123456789012:trail/my-trail", results[0].ARN)
	assert.Equal(t, "my-trail-bucket", results[0].S3Bucket)
}

func TestCloudTrailFetchAll_Empty(t *testing.T) {
	mockClient := new(mockCloudTrailClient)
	mockClient.On("DescribeTrails", mock.Anything, mock.Anything).Return(&cloudtrail.DescribeTrailsOutput{
		TrailList: []types.Trail{},
	}, nil)

	fetcher := NewCloudTrailFetcher(mockClient)
	results, err := fetcher.FetchAll(context.Background())

	assert.NoError(t, err)
	assert.Empty(t, results)
}
