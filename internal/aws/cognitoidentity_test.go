package aws_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentity/types"
	"github.com/stretchr/testify/assert"
)

type mockCognitoIdentityClient struct {
	pools []types.IdentityPoolShortDescription
	err   error
}

func (m *mockCognitoIdentityClient) ListIdentityPools(ctx context.Context, params *cognitoidentity.ListIdentityPoolsInput, optFns ...func(*cognitoidentity.Options)) (*cognitoidentity.ListIdentityPoolsOutput, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &cognitoidentity.ListIdentityPoolsOutput{
		IdentityPools: m.pools,
	}, nil
}

func TestCognitoIdentityFetchAll_ReturnsPools(t *testing.T) {
	mock := &mockCognitoIdentityClient{
		pools: []types.IdentityPoolShortDescription{
			{IdentityPoolId: aws.String("us-east-1:abc-123"), IdentityPoolName: aws.String("my-pool")},
			{IdentityPoolId: aws.String("us-east-1:def-456"), IdentityPoolName: aws.String("other-pool")},
		},
	}

	resources, err := fetchCognitoIdentityPoolsWithClient(mock)
	assert.NoError(t, err)
	assert.Len(t, resources, 2)
	assert.Equal(t, "us-east-1:abc-123", resources[0].ID)
	assert.Equal(t, "my-pool", resources[0].Name)
}

func TestCognitoIdentityFetchAll_Empty(t *testing.T) {
	mock := &mockCognitoIdentityClient{
		pools: []types.IdentityPoolShortDescription{},
	}

	resources, err := fetchCognitoIdentityPoolsWithClient(mock)
	assert.NoError(t, err)
	assert.Empty(t, resources)
}
