package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
	"github.com/aws/aws-sdk-go-v2/service/kinesis/types"
)

type KinesisStream struct {
	Name   string
	ARN    string
	Status string
}

type kinesisClient interface {
	ListStreams(ctx context.Context, params *kinesis.ListStreamsInput, optFns ...func(*kinesis.Options)) (*kinesis.ListStreamsOutput, error)
	DescribeStreamSummary(ctx context.Context, params *kinesis.DescribeStreamSummaryInput, optFns ...func(*kinesis.Options)) (*kinesis.DescribeStreamSummaryOutput, error)
}

type KinesisFetcher struct {
	client kinesisClient
}

func NewKinesisFetcher(client kinesisClient) *KinesisFetcher {
	return &KinesisFetcher{client: client}
}

func (f *KinesisFetcher) FetchAll(ctx context.Context) ([]KinesisStream, error) {
	var streams []KinesisStream
	var nextToken *string

	for {
		resp, err := f.client.ListStreams(ctx, &kinesis.ListStreamsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, summary := range resp.StreamSummaries {
			streams = append(streams, mapKinesisStream(summary))
		}

		if !resp.HasMoreStreams || resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return streams, nil
}

func mapKinesisStream(s types.StreamSummary) KinesisStream {
	return KinesisStream{
		Name:   aws.ToString(s.StreamName),
		ARN:    aws.ToString(s.StreamARN),
		Status: string(s.StreamStatus),
	}
}
