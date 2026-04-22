package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/config"
)

type CognitoIdentityPool struct {
	ID         string
	Name       string
	AllowUnauthenticated bool
}

type CognitoIdentityClient interface {
	ListIdentityPools(ctx context.Context, params *cognitoidentity.ListIdentityPoolsInput, optFns ...func(*cognitoidentity.Options)) (*cognitoidentity.ListIdentityPoolsOutput, error)
}

type CognitoIdentityFetcher struct {
	client CognitoIdentityClient
}

func NewCognitoIdentityFetcher(ctx context.Context, region string) (*CognitoIdentityFetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &CognitoIdentityFetcher{
		client: cognitoidentity.NewFromConfig(cfg),
	}, nil
}

func (f *CognitoIdentityFetcher) FetchAll(ctx context.Context) ([]CognitoIdentityPool, error) {
	var pools []CognitoIdentityPool
	var nextToken *string

	for {
		out, err := f.client.ListIdentityPools(ctx, &cognitoidentity.ListIdentityPoolsInput{
			MaxResults: aws.Int32(60),
			NextToken:  nextToken,
		})
		if err != nil {
			return nil, err
		}
		for _, p := range out.IdentityPools {
			pools = append(pools, mapIdentityPool(p))
		}
		if out.NextToken == nil {
			break
		}
		nextToken = out.NextToken
	}
	return pools, nil
}

func mapIdentityPool(p cognitoidentity.IdentityPoolShortDescription) CognitoIdentityPool {
	return CognitoIdentityPool{
		ID:   aws.ToString(p.IdentityPoolId),
		Name: aws.ToString(p.IdentityPoolName),
	}
}
