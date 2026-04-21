package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

type mockACMClient struct {
	output *acm.ListCertificatesOutput
	err    error
}

func (m *mockACMClient) ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error) {
	return m.output, m.err
}

func TestACMFetchAll_ReturnsCertificates(t *testing.T) {
	mock := &mockACMClient{
		output: &acm.ListCertificatesOutput{
			CertificateSummaryList: []types.CertificateSummary{
				{
					CertificateArn: aws.String("arn:aws:acm:us-east-1:123456789012:certificate/abc-123"),
					DomainName:     aws.String("example.com"),
					Status:         types.CertificateStatusIssued,
				},
			},
		},
	}
	fetcher := NewACMFetcher(mock)
	certs, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(certs) != 1 {
		t.Fatalf("expected 1 certificate, got %d", len(certs))
	}
	if certs[0].Domain != "example.com" {
		t.Errorf("expected domain example.com, got %s", certs[0].Domain)
	}
	if certs[0].Status != "ISSUED" {
		t.Errorf("expected status ISSUED, got %s", certs[0].Status)
	}
}

func TestACMFetchAll_Empty(t *testing.T) {
	mock := &mockACMClient{
		output: &acm.ListCertificatesOutput{},
	}
	fetcher := NewACMFetcher(mock)
	certs, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(certs) != 0 {
		t.Errorf("expected 0 certificates, got %d", len(certs))
	}
}
