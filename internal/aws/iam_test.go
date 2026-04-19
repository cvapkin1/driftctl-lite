package aws

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
)

type mockIAMClient struct {
	roles []types.Role
}

func (m *mockIAMClient) ListRoles(_ context.Context, _ *iam.ListRolesInput, _ ...func(*iam.Options)) (*iam.ListRolesOutput, error) {
	return &iam.ListRolesOutput{Roles: m.roles, IsTruncated: false}, nil
}

func TestIAMFetchAll_ReturnRoles(t *testing.T) {
	mock := &mockIAMClient{
		roles: []types.Role{
			{RoleId: aws.String("AROA123"), RoleName: aws.String("MyRole"), Arn: aws.String("arn:aws:iam::123456789012:role/MyRole")},
		},
	}
	fetcher := NewIAMFetcher(mock)
	roles, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 1 {
		t.Fatalf("expected 1 role, got %d", len(roles))
	}
	if roles[0].Name != "MyRole" {
		t.Errorf("expected role name MyRole, got %s", roles[0].Name)
	}
}

func TestIAMFetchAll_Empty(t *testing.T) {
	mock := &mockIAMClient{roles: []types.Role{}}
	fetcher := NewIAMFetcher(mock)
	roles, err := fetcher.FetchAll(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(roles) != 0 {
		t.Errorf("expected 0 roles, got %d", len(roles))
	}
}
