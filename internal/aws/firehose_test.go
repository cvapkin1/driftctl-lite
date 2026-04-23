package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFirehoseClient struct {
	listOut  *firehose.ListDeliveryStreamsOutput
	descOut  *firehose.DescribeDeliveryStreamOutput
	listErr  error
	descErr  error
}

func (m *mockFirehoseClient) ListDeliveryStreams(ctx context.Context, params *firehose.ListDeliveryStreamsInput, optFns ...func(*firehose.Options)) (*firehose.ListDeliveryStreamsOutput, error) {
	return m.listOut, m.listErr
}

func (m *mockFirehoseClient) DescribeDeliveryStream(ctx context.Context, params *firehose.DescribeDeliveryStreamInput, optFns ...func(*firehose.Options)) (*firehose.DescribeDeliveryStreamOutput, error) {
	return m.descOut, m.descErr
}

func TestFirehoseFetchAll_ReturnsStreams(t *testing.T) {
	mock := &mockFirehoseClient{
		listOut: &firehose.ListDeliveryStreamsOutput{
			DeliveryStreamNames: []string{"my-stream"},
			HasMoreDeliveryStreams: false,
		},
		descOut: &firehose.DescribeDeliveryStreamOutput{
			DeliveryStreamDescription: &types.DeliveryStreamDescription{
				DeliveryStreamName:   aws.String("my-stream"),
				DeliveryStreamARN:    aws.String("arn:aws:firehose:us-east-1:123456789012:deliverystream/my-stream"),
				DeliveryStreamStatus: types.DeliveryStreamStatusActive,
			},
		},
	}

	fetcher := NewFirehoseFetcher(mock)
	streams, err := fetcher.FetchAll(context.Background())

	require.NoError(t, err)
	require.Len(t, streams, 1)
	assert.Equal(t, "my-stream", streams[0].Name)
	assert.Equal(t, "arn:aws:firehose:us-east-1:123456789012:deliverystream/my-stream", streams[0].ARN)
	assert.Equal(t, "ACTIVE", streams[0].Status)
}

func TestFirehoseFetchAll_Empty(t *testing.T) {
	mock := &mockFirehoseClient{
		listOut: &firehose.ListDeliveryStreamsOutput{
			DeliveryStreamNames:   []string{},
			HasMoreDeliveryStreams: false,
		},
	}

	fetcher := NewFirehoseFetcher(mock)
	streams, err := fetcher.FetchAll(context.Background())

	require.NoError(t, err)
	assert.Empty(t, streams)
}
