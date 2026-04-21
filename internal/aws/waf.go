package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

type WAFWebACL struct {
	ID          string
	Name        string
	ARN         string
	Scope       string
	Description string
}

type WAFFetcher struct {
	client *wafv2.Client
}

type wafv2API interface {
	ListWebACLs(ctx context.Context, params *wafv2.ListWebACLsInput, optFns ...func(*wafv2.Options)) (*wafv2.ListWebACLsOutput, error)
}

func NewWAFFetcher(client *wafv2.Client) *WAFFetcher {
	return &WAFFetcher{client: client}
}

func (f *WAFFetcher) FetchAll(ctx context.Context) ([]WAFWebACL, error) {
	var results []WAFWebACL

	for _, scope := range []types.Scope{types.ScopeRegional, types.ScopeCloudfront} {
		var nextMarker *string
		for {
			out, err := f.client.ListWebACLs(ctx, &wafv2.ListWebACLsInput{
				Scope:      scope,
				NextMarker: nextMarker,
				Limit:      aws.Int32(100),
			})
			if err != nil {
				return nil, err
			}
			for _, acl := range out.WebACLs {
				results = append(results, mapWebACL(acl, string(scope)))
			}
			if out.NextMarker == nil {
				break
			}
			nextMarker = out.NextMarker
		}
	}

	return results, nil
}

func mapWebACL(acl types.WebACLSummary, scope string) WAFWebACL {
	return WAFWebACL{
		ID:          aws.ToString(acl.Id),
		Name:        aws.ToString(acl.Name),
		ARN:         aws.ToString(acl.ARN),
		Scope:       scope,
		Description: aws.ToString(acl.Description),
	}
}
