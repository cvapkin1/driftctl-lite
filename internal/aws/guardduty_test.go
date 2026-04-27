package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	"github.com/aws/aws-sdk-go-v2/service/guardduty/types"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func TestMapGuardDutyDetector_Fields(t *testing.T) {
	out := &guardduty.GetDetectorOutput{
		Status:      types.DetectorStatus("ENABLED"),
		ServiceRole: aws.String("arn:aws:iam::123456789012:role/aws-service-role/guardduty.amazonaws.com/AWSServiceRoleForAmazonGuardDuty"),
	}
	d := mapGuardDutyDetector("abc123", out)

	if d.ID != "abc123" {
		t.Errorf("expected ID abc123, got %s", d.ID)
	}
	if d.Status != "ENABLED" {
		t.Errorf("expected status ENABLED, got %s", d.Status)
	}
	if d.ARN == "" {
		t.Error("expected non-empty ARN")
	}
}

func TestMapGuardDutyDetector_EmptyFields(t *testing.T) {
	out := &guardduty.GetDetectorOutput{}
	d := mapGuardDutyDetector("xyz", out)

	if d.ID != "xyz" {
		t.Errorf("expected ID xyz, got %s", d.ID)
	}
	if d.Status != "" {
		t.Errorf("expected empty status, got %s", d.Status)
	}
}
