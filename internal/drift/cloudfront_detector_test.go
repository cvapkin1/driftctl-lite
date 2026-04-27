package drift_test

import (
	"testing"

	"github.com/your-org/driftctl-lite/internal/drift"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

func cloudfrontStateResource(id, domainName, status string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_cloudfront_distribution",
		Name: id,
		Attributes: map[string]string{
			"id":          id,
			"domain_name": domainName,
			"status":      status,
		},
	}
}

func TestDetectCloudFrontDrift_Missing(t *testing.T) {
	state := []tfstate.Resource{
		cloudfrontStateResource("EDFDVBD6EXAMPLE", "d111111abcdef8.cloudfront.net", "Deployed"),
	}
	live := []map[string]string{}

	results := drift.DetectCloudFrontDrift(state, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].Status != "missing" {
		t.Errorf("expected status 'missing', got %s", results[0].Status)
	}
}

func TestDetectCloudFrontDrift_StatusMismatch(t *testing.T) {
	state := []tfstate.Resource{
		cloudfrontStateResource("EDFDVBD6EXAMPLE", "d111111abcdef8.cloudfront.net", "Deployed"),
	}
	live := []map[string]string{
		{
			"id":          "EDFDVBD6EXAMPLE",
			"domain_name": "d111111abcdef8.cloudfront.net",
			"status":      "InProgress",
		},
	}

	results := drift.DetectCloudFrontDrift(state, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].Status != "changed" {
		t.Errorf("expected status 'changed', got %s", results[0].Status)
	}
}

func TestDetectCloudFrontDrift_DomainMismatch(t *testing.T) {
	state := []tfstate.Resource{
		cloudfrontStateResource("EDFDVBD6EXAMPLE", "d111111abcdef8.cloudfront.net", "Deployed"),
	}
	live := []map[string]string{
		{
			"id":          "EDFDVBD6EXAMPLE",
			"domain_name": "dXXXXXXXXXXXX.cloudfront.net",
			"status":      "Deployed",
		},
	}

	results := drift.DetectCloudFrontDrift(state, live)
	if len(results) != 1 {
		t.Fatalf("expected 1 drift result, got %d", len(results))
	}
	if results[0].Status != "changed" {
		t.Errorf("expected status 'changed', got %s", results[0].Status)
	}
}

func TestDetectCloudFrontDrift_NoDrift(t *testing.T) {
	state := []tfstate.Resource{
		cloudfrontStateResource("EDFDVBD6EXAMPLE", "d111111abcdef8.cloudfront.net", "Deployed"),
	}
	live := []map[string]string{
		{
			"id":          "EDFDVBD6EXAMPLE",
			"domain_name": "d111111abcdef8.cloudfront.net",
			"status":      "Deployed",
		},
	}

	results := drift.DetectCloudFrontDrift(state, live)
	if len(results) != 0 {
		t.Errorf("expected no drift, got %d results", len(results))
	}
}
