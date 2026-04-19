package drift_test

import (
	"testing"

	"github.com/example/driftctl-lite/internal/aws"
	"github.com/example/driftctl-lite/internal/drift"
	"github.com/example/driftctl-lite/internal/tfstate"
)

func stateResource(id string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_instance",
		Name: "web",
		Attributes: map[string]interface{}{"id": id},
	}
}

func TestDetectEC2Drift_Missing(t *testing.T) {
	res := []tfstate.Resource{stateResource("i-missing")}
	results := drift.DetectEC2Drift(res, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift, got %d", len(results))
	}
	if results[0].Reason != "resource not found in AWS" {
		t.Errorf("unexpected reason: %s", results[0].Reason)
	}
}

func TestDetectEC2Drift_Terminated(t *testing.T) {
	res := []tfstate.Resource{stateResource("i-term")}
	live := []aws.EC2Instance{{ID: "i-term", State: "terminated"}}
	results := drift.DetectEC2Drift(res, live)
	if len(results) != 1 || results[0].Reason != "instance is terminated" {
		t.Errorf("expected terminated drift, got %+v", results)
	}
}

func TestDetectEC2Drift_NoDrift(t *testing.T) {
	res := []tfstate.Resource{stateResource("i-ok")}
	live := []aws.EC2Instance{{ID: "i-ok", State: "running"}}
	results := drift.DetectEC2Drift(res, live)
	if len(results) != 0 {
		t.Errorf("expected no drift, got %+v", results)
	}
}
