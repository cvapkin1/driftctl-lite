package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

type mockRoute53Client struct {
	zones []types.HostedZone
}

func (m *mockRoute53Client) ListHostedZones(ctx context.Context, params *route53.ListHostedZonesInput, optFns ...func(*route53.Options)) (*route53.ListHostedZonesOutput, error) {
	return &route53.ListHostedZonesOutput{HostedZones: m.zones}, nil
}

func TestRoute53FetchAll_ReturnsZones(t *testing.T) {
	mock := &mockRoute53Client{
		zones: []types.HostedZone{
			{Id: aws.String("/hostedzone/Z123"), Name: aws.String("example.com.")},
			{Id: aws.String("/hostedzone/Z456"), Name: aws.String("test.io.")},
		},
	}
	fetcher := NewRoute53Fetcher(mock)
	zones, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(zones) != 2 {
		t.Fatalf("expected 2 zones, got %d", len(zones))
	}
	if zones[0].ID != "/hostedzone/Z123" || zones[0].Name != "example.com." {
		t.Errorf("unexpected zone: %+v", zones[0])
	}
}

func TestRoute53FetchAll_Empty(t *testing.T) {
	mock := &mockRoute53Client{zones: []types.HostedZone{}}
	fetcher := NewRoute53Fetcher(mock)
	zones, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(zones) != 0 {
		t.Errorf("expected 0 zones, got %d", len(zones))
	}
}
