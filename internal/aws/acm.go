package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
)

type ACMCertificate struct {
	ARN    string
	Domain string
	Status string
}

type ACMClient interface {
	ListCertificates(ctx context.Context, params *acm.ListCertificatesInput, optFns ...func(*acm.Options)) (*acm.ListCertificatesOutput, error)
}

type ACMFetcher struct {
	client ACMClient
}

func NewACMFetcher(client ACMClient) *ACMFetcher {
	return &ACMFetcher{client: client}
}

func (f *ACMFetcher) FetchAll(ctx context.Context) ([]ACMCertificate, error) {
	var certs []ACMCertificate
	var nextToken *string

	for {
		out, err := f.client.ListCertificates(ctx, &acm.ListCertificatesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, c := range out.CertificateSummaryList {
			certs = append(certs, mapCertificate(c))
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return certs, nil
}

func mapCertificate(c types.CertificateSummary) ACMCertificate {
	return ACMCertificate{
		ARN:    aws.ToString(c.CertificateArn),
		Domain: aws.ToString(c.DomainName),
		Status: string(c.Status),
	}
}
