package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type IAMRole struct {
	ID   string
	Name string
	ARN  string
	Tags map[string]string
}

type interface {
	ListRoles(ctx context.Context, params *iam.ListRolesInput, optFns ...func(*iam.Options)) (*iam.ListRolesOutput, error)
}

type IAMFetcher struct {
	client IAMClient
}

func NewIAMFetcher(client IAMClient) *IAMFetcher {
	return &IAMFetcher{client: client}
}

func (f *IAMFetcher) FetchAll(ctx context.Context) ([]IAMRole, error) {
	var roles []IAMRole
	var marker *string

	for {
		out, err := f.client.ListRoles(ctx, &iam.ListRolesInput{Marker: marker})
		if err != nil {
			return nil, err
		}
		for _, r := range out.Roles {
			roles = append(roles, IAMRole{
				ID:   aws.ToString(r.RoleId),
				Name: aws.ToString(r.RoleName),
				ARN:  aws.ToString(r.Arn),
				Tags: map[string]string{},
			})
		}
		if !out.IsTruncated {
			break
		}
		marker = out.Marker
	}
	return roles, nil
}
