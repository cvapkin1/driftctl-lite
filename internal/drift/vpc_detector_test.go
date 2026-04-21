package drift

import (
	"testing"

	"github.com/owner/driftctl-lite/internal/aws"
)

func vpcStateResource(id, cidr string) map[string]interface{} {
	return map[string]interface{}{
		"id":         id,
		"cidr_block": cidr,
	}
}

func TestDetectVPCDrift_Missing(t *testing.T) {
	state := []map[string]interface{}{vpcStateResource("vpc-111", "10.0.0.0/16")}
	live := []aws.VPCResource{}

	results := DetectVPCDrift(state, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].DriftType != "missing" {
		t.Errorf("expected drift type 'missing', got %s", results[0].DriftType)
	}
}

func TestDetectVPCDrift_StateMismatch(t *testing.T) {
	state := []map[string]interface{}{vpcStateResource("vpc-222", "172.16.0.0/12")}
	live := []aws.VPCResource{
		{ID: "vpc-222", CIDR: "172.16.0.0/12", State: "pending"},
	}

	results := DetectVPCDrift(state, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].DriftType != "state_mismatch" {
		t.Errorf("expected state_mismatch, got %s", results[0].DriftType)
	}
}

func TestDetectVPCDrift_CIDRMismatch(t *testing.T) {
	state := []map[string]interface{}{vpcStateResource("vpc-333", "10.1.0.0/16")}
	live := []aws.VPCResource{
		{ID: "vpc-333", CIDR: "10.2.0.0/16", State: "available"},
	}

	results := DetectVPCDrift(state, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].DriftType != "cidr_mismatch" {
		t.Errorf("expected cidr_mismatch, got %s", results[0].DriftType)
	}
}

func TestDetectVPCDrift_NoDrift(t *testing.T) {
	state := []map[string]interface{}{vpcStateResource("vpc-444", "192.168.0.0/24")}
	live := []aws.VPCResource{
		{ID: "vpc-444", CIDR: "192.168.0.0/24", State: "available"},
	}

	results := DetectVPCDrift(state, live)
	if len(results) != 0 {
		t.Errorf("expected no drift, got %d results", len(results))
	}
}
