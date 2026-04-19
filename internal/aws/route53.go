package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type Route53Zone struct {
	ID   string
	Name string
}

type route53Client interface {
	ListHostedZones(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error)
}

type Route53Fetcher struct {
	client route53Client
}

func NewRoute53Fetcher(client route53Client) *Route53Fetcher {
	return &Route53Fetcher{client: client}
}

func (f *Route53Fetcher) FetchAll(ctx context.Context) ([]Route53Zone, error) {
	out, err := f.client.ListHostedZones(ctx, &route53.ListHostedZonesInput{})
	if err != nil {
		return nil, err
	}

	zones := make([]Route53Zone, 0, len(out.HostedZones))
	for _, z := range out.HostedZones {
		zones = append(zones, mapZone(z))
	}
	return zones, nil
}

func mapZone(z types.HostedZone) Route53Zone {
	return Route53Zone{
		ID:   aws.ToString(z.Id),
		Name: aws.ToString(z.Name),
	}
}
