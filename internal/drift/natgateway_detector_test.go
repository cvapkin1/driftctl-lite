package drift

import (
	"testing"

	"github.com/jonhadfield/driftctl-lite/internal/aws"
)

func natgwStateResource(id, subnetID string) Resource {
	return Resource{
		Type: "aws_nat_gateway",
		Attributes: map[string]interface{}{
			"id":        id,
			"subnet_id": subnetID,
		},
	}
}

func TestDetectNATGatewayDrift_Missing(t *testing.T) {
	state := []Resource{natgwStateResource("nat-001", "subnet-aaa")}
	live := []aws.NATGateway{}

	results := DetectNATGatewayDrift(state, live)

	if len(results) != 1 || results[0].DriftType != "missing" {
		t.Errorf("expected missing drift, got %+v", results)
	}
}

func TestDetectNATGatewayDrift_Deleted(t *testing.T) {
	state := []Resource{natgwStateResource("nat-002", "subnet-bbb")}
	live := []aws.NATGateway{
		{ID: "nat-002", State: "deleted", SubnetID: "subnet-bbb"},
	}

	results := DetectNATGatewayDrift(state, live)

	if len(results) != 1 || results[0].DriftType != "deleted" {
		t.Errorf("expected deleted drift, got %+v", results)
	}
}

func TestDetectNATGatewayDrift_SubnetMismatch(t *testing.T) {
	state := []Resource{natgwStateResource("nat-003", "subnet-expected")}
	live := []aws.NATGateway{
		{ID: "nat-003", State: "available", SubnetID: "subnet-actual"},
	}

	results := DetectNATGatewayDrift(state, live)

	if len(results) != 1 || results[0].DriftType != "subnet_mismatch" {
		t.Errorf("expected subnet_mismatch drift, got %+v", results)
	}
}

func TestDetectNATGatewayDrift_NoDrift(t *testing.T) {
	state := []Resource{natgwStateResource("nat-004", "subnet-ccc")}
	live := []aws.NATGateway{
		{ID: "nat-004", State: "available", SubnetID: "subnet-ccc"},
	}

	results := DetectNATGatewayDrift(state, live)

	if len(results) != 0 {
		t.Errorf("expected no drift, got %+v", results)
	}
}
