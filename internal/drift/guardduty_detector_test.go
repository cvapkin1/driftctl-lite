package drift

import (
	"testing"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

func guarddutyStateResource(id string, enable bool) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_guardduty_detector",
		Name: "main",
		Attributes: map[string]interface{}{
			"id":     id,
			"enable": enable,
		},
	}
}

func TestDetectGuardDutyDrift_Missing(t *testing.T) {
	res := guarddutyStateResource("det-001", true)
	results := DetectGuardDutyDrift([]tfstate.Resource{res}, []aws.GuardDutyDetector{})

	if len(results) != 1 {
		t.Fatalf("expected 1 drift, got %d", len(results))
	}
	if results[0].DriftType != "missing" {
		t.Errorf("expected missing, got %s", results[0].DriftType)
	}
}

func TestDetectGuardDutyDrift_Disabled(t *testing.T) {
	res := guarddutyStateResource("det-001", true)
	live := []aws.GuardDutyDetector{{ID: "det-001", Status: "DISABLED"}}
	results := DetectGuardDutyDrift([]tfstate.Resource{res}, live)

	if len(results) != 1 {
		t.Fatalf("expected 1 drift, got %d", len(results))
	}
	if results[0].DriftType != "disabled" {
		t.Errorf("expected disabled, got %s", results[0].DriftType)
	}
}

func TestDetectGuardDutyDrift_StatusMismatch(t *testing.T) {
	res := guarddutyStateResource("det-002", true)
	live := []aws.GuardDutyDetector{{ID: "det-002", Status: "PENDING"}}
	results := DetectGuardDutyDrift([]tfstate.Resource{res}, live)

	if len(results) != 1 {
		t.Fatalf("expected 1 drift, got %d", len(results))
	}
	if results[0].DriftType != "status_mismatch" {
		t.Errorf("expected status_mismatch, got %s", results[0].DriftType)
	}
}

func TestDetectGuardDutyDrift_NoDrift(t *testing.T) {
	res := guarddutyStateResource("det-003", true)
	live := []aws.GuardDutyDetector{{ID: "det-003", Status: "ENABLED"}}
	results := DetectGuardDutyDrift([]tfstate.Resource{res}, live)

	if len(results) != 0 {
		t.Errorf("expected no drift, got %d", len(results))
	}
}
