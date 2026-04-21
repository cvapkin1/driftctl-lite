package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/wafv2/types"
)

func TestMapWebACL_Fields(t *testing.T) {
	acl := types.WebACLSummary{
		Id:          aws.String("abc-123"),
		Name:        aws.String("my-web-acl"),
		ARN:         aws.String("arn:aws:wafv2:us-east-1:123456789012:regional/webacl/my-web-acl/abc-123"),
		Description: aws.String("blocks bad traffic"),
	}

	result := mapWebACL(acl, "REGIONAL")

	if result.ID != "abc-123" {
		t.Errorf("expected ID abc-123, got %s", result.ID)
	}
	if result.Name != "my-web-acl" {
		t.Errorf("expected Name my-web-acl, got %s", result.Name)
	}
	if result.ARN != "arn:aws:wafv2:us-east-1:123456789012:regional/webacl/my-web-acl/abc-123" {
		t.Errorf("unexpected ARN: %s", result.ARN)
	}
	if result.Scope != "REGIONAL" {
		t.Errorf("expected Scope REGIONAL, got %s", result.Scope)
	}
	if result.Description != "blocks bad traffic" {
		t.Errorf("expected description, got %s", result.Description)
	}
}

func TestMapWebACL_EmptyDescription(t *testing.T) {
	acl := types.WebACLSummary{
		Id:   aws.String("xyz-456"),
		Name: aws.String("minimal-acl"),
		ARN:  aws.String("arn:aws:wafv2:us-east-1:123456789012:regional/webacl/minimal-acl/xyz-456"),
	}

	result := mapWebACL(acl, "CLOUDFRONT")

	if result.Description != "" {
		t.Errorf("expected empty description, got %s", result.Description)
	}
	if result.Scope != "CLOUDFRONT" {
		t.Errorf("expected Scope CLOUDFRONT, got %s", result.Scope)
	}
}
