package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// EC2Instance represents a live EC2 instance fetched from AWS.
type EC2Instance struct {
	ID           string
	InstanceType string
	State        string
	Tags         map[string]string
}

// EC2Fetcher fetches live EC2 instances from AWS.
type EC2Fetcher struct {
	client *ec2.Client
}

// NewEC2Fetcher creates a new EC2Fetcher using the provided AWS config.
func NewEC2Fetcher(cfg aws.Config) *EC2Fetcher {
	return &EC2Fetcher{client: ec2.NewFromConfig(cfg)}
}

// FetchAll returns all EC2 instances visible to the caller.
func (f *EC2Fetcher) FetchAll(ctx context.Context) ([]EC2Instance, error) {
	var instances []EC2Instance
	paginator := ec2.NewDescribeInstancesPaginator(f.client, &ec2.DescribeInstancesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("describing instances: %w", err)
		}
		for _, r := range page.Reservations {
			for _, i := range r.Instances {
				instances = append(instances, mapInstance(i))
			}
		}
	}
	return instances, nil
}

func mapInstance(i types.Instance) EC2Instance {
	tags := make(map[string]string, len(i.Tags))
	for _, t := range i.Tags {
		if t.Key != nil && t.Value != nil {
			tags[*t.Key] = *t.Value
		}
	}
	return EC2Instance{
		ID:           aws.ToString(i.InstanceId),
		InstanceType: string(i.InstanceType),
		State:        string(i.State.Name),
		Tags:         tags,
	}
}
