package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
	"github.com/aws/aws-sdk-go-v2/service/firehose/types"
)

type FirehoseDeliveryStream struct {
	Name   string
	ARN    string
	Status string
}

type firehoseClient interface {
	ListDeliveryStreams(ctx context.Context, params *firehose.ListDeliveryStreamsInput, optFns ...func(*firehose.Options)) (*firehose.ListDeliveryStreamsOutput, error)
	DescribeDeliveryStream(ctx context.Context, params *firehose.DescribeDeliveryStreamInput, optFns ...func(*firehose.Options)) (*firehose.DescribeDeliveryStreamOutput, error)
}

type FirehoseFetcher struct {
	client firehoseClient
}

func NewFirehoseFetcher(client firehoseClient) *FirehoseFetcher {
	return &FirehoseFetcher{client: client}
}

func (f *FirehoseFetcher) FetchAll(ctx context.Context) ([]FirehoseDeliveryStream, error) {
	var streams []FirehoseDeliveryStream
	var exclusiveStart *string

	for {
		listOut, err := f.client.ListDeliveryStreams(ctx, &firehose.ListDeliveryStreamsInput{
			ExclusiveStartDeliveryStreamName: exclusiveStart,
		})
		if err != nil {
			return nil, err
		}

		for _, name := range listOut.DeliveryStreamNames {
			descOut, err := f.client.DescribeDeliveryStream(ctx, &firehose.DescribeDeliveryStreamInput{
				DeliveryStreamName: aws.String(name),
			})
			if err != nil {
				return nil, err
			}
			streams = append(streams, mapDeliveryStream(descOut.DeliveryStreamDescription))
		}

		if !listOut.HasMoreDeliveryStreams {
			break
		}
		if len(listOut.DeliveryStreamNames) > 0 {
			exclusiveStart = aws.String(listOut.DeliveryStreamNames[len(listOut.DeliveryStreamNames)-1])
		}
	}

	return streams, nil
}

func mapDeliveryStream(d *types.DeliveryStreamDescription) FirehoseDeliveryStream {
	if d == nil {
		return FirehoseDeliveryStream{}
	}
	return FirehoseDeliveryStream{
		Name:   aws.ToString(d.DeliveryStreamName),
		ARN:    aws.ToString(d.DeliveryStreamARN),
		Status: string(d.DeliveryStreamStatus),
	}
}
