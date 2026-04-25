package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	"github.com/aws/aws-sdk-go-v2/service/lightsail/types"
)

type LightsailInstance struct {
	ID    string
	Name  string
	State string
	ARN   string
}

type lightsailClient interface {
	GetInstances(ctx context.Context, params *lightsail.GetInstancesInput, optFns ...func(*lightsail.Options)) (*lightsail.GetInstancesOutput, error)
}

type LightsailFetcher struct {
	client lightsailClient
}

func NewLightsailFetcher(client lightsailClient) *LightsailFetcher {
	return &LightsailFetcher{client: client}
}

func (f *LightsailFetcher) FetchAll(ctx context.Context) ([]LightsailInstance, error) {
	var instances []LightsailInstance
	var pageToken *string

	for {
		out, err := f.client.GetInstances(ctx, &lightsail.GetInstancesInput{
			PageToken: pageToken,
		})
		if err != nil {
			return nil, err
		}

		for _, inst := range out.Instances {
			instances = append(instances, mapLightsailInstance(inst))
		}

		if out.NextPageToken == nil {
			break
		}
		pageToken = out.NextPageToken
	}

	return instances, nil
}

func mapLightsailInstance(inst types.Instance) LightsailInstance {
	state := ""
	if inst.State != nil {
		state = aws.ToString(inst.State.Name)
	}
	return LightsailInstance{
		ID:    aws.ToString(inst.SupportCode),
		Name:  aws.ToString(inst.Name),
		State: state,
		ARN:   aws.ToString(inst.Arn),
	}
}
