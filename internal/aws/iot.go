package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

type IoTThing struct {
	ThingName string
	ThingType string
	ARN       string
}

type IoTFetcher struct {
	client *iot.Client
}

func NewIoTFetcher(client *iot.Client) *IoTFetcher {
	return &IoTFetcher{client: client}
}

func (f *IoTFetcher) FetchAll(ctx context.Context) ([]IoTThing, error) {
	var things []IoTThing
	var nextToken *string

	for {
		resp, err := f.client.ListThings(ctx, &iot.ListThingsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, t := range resp.Things {
			things = append(things, mapIoTThing(t))
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return things, nil
}

func mapIoTThing(t iot.ThingAttribute) IoTThing {
	name := aws.ToString(t.ThingName)
	thingType := aws.ToString(t.ThingTypeName)
	arn := aws.ToString(t.ThingArn)
	return IoTThing{
		ThingName: name,
		ThingType: thingType,
		ARN:       arn,
	}
}
