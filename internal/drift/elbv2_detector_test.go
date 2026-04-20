package drift_test

import (
	"testing"

	"github.com/snyk/driftctl-lite/internal/drift"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

func elbv2StateResource(id, name, lbType string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_lb",
		Name: name,
		Attributes: map[string]interface{}{
			"id":                id,
			"name":              name,
			"load_balancer_type": lbType,
			"dns_name":          name + ".us-east-1.elb.amazonaws.com",
		},
	}
}

func TestDetectELBv2Drift_Missing(t *testing.T) {
	stateResources := []tfstate.Resource{
		elbv2StateResource("arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-lb/abc123", "my-lb", "application"),
	}

	// No live load balancers
	live := map[string]string{}

	results := drift.DetectELBv2Drift(stateResources, live)

	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].ResourceID != "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-lb/abc123" {
		t.Errorf("unexpected resource ID: %s", results[0].ResourceID)
	}
	if results[0].Status != "missing" {
		t.Errorf("expected status 'missing', got '%s'", results[0].Status)
	}
}

func TestDetectELBv2Drift_DNSMismatch(t *testing.T) {
	arn := "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-lb/abc123"
	stateResources := []tfstate.Resource{
		elbv2StateResource(arn, "my-lb", "application"),
	}

	// Live LB has a different DNS name
	live := map[string]string{
		arn: "different-dns.us-east-1.elb.amazonaws.com",
	}

	results := drift.DetectELBv2Drift(stateResources, live)

	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].Status != "dns_mismatch" {
		t.Errorf("expected status 'dns_mismatch', got '%s'", results[0].Status)
	}
}

func TestDetectELBv2Drift_NoDrift(t *testing.T) {
	arn := "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/my-lb/abc123"
	stateResources := []tfstate.Resource{
		elbv2StateResource(arn, "my-lb", "application"),
	}

	// Live LB matches state
	live := map[string]string{
		arn: "my-lb.us-east-1.elb.amazonaws.com",
	}

	results := drift.DetectELBv2Drift(stateResources, live)

	if len(results) != 0 {
		t.Errorf("expected 0 drift results, got %d", len(results))
	}
}

func TestDetectELBv2Drift_MultipleResources(t *testing.T) {
	arn1 := "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/lb-one/aaa"
	arn2 := "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/lb-two/bbb"

	stateResources := []tfstate.Resource{
		elbv2StateResource(arn1, "lb-one", "application"),
		elbv2StateResource(arn2, "lb-two", "application"),
	}

	// Only first LB exists live
	live := map[string]string{
		arn1: "lb-one.us-east-1.elb.amazonaws.com",
	}

	results := drift.DetectELBv2Drift(stateResources, live)

	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].ResourceID != arn2 {
		t.Errorf("expected missing resource %s, got %s", arn2, results[0].ResourceID)
	}
	if results[0].Status != "missing" {
		t.Errorf("expected status 'missing', got '%s'", results[0].Status)
	}
}
