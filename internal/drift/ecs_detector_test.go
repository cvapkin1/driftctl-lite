package drift

import (
	"testing"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

func ecsStateResource(arn, name string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_ecs_cluster",
		Name: name,
		Attributes: map[string]interface{}{
			"arn":  arn,
			"name": name,
		},
	}
}

func TestDetectECSDrift_Missing(t *testing.T) {
	state := []tfstate.Resource{ecsStateResource("arn:aws:ecs:::cluster/gone", "gone")}
	live := []aws.ECSCluster{}

	results := DetectECSDrift(state, live)

	if len(results) != 1 || results[0].DriftType != "missing" {
		t.Errorf("expected 1 missing drift, got %+v", results)
	}
}

func TestDetectECSDrift_Inactive(t *testing.T) {
	arn := "arn:aws:ecs:::cluster/inactive"
	state := []tfstate.Resource{ecsStateResource(arn, "inactive")}
	live := []aws.ECSCluster{{ARN: arn, Name: "inactive", Status: "INACTIVE"}}

	results := DetectECSDrift(state, live)

	if len(results) != 1 || results[0].DriftType != "deleted" {
		t.Errorf("expected 1 deleted drift, got %+v", results)
	}
}

func TestDetectECSDrift_NameMismatch(t *testing.T) {
	arn := "arn:aws:ecs:::cluster/renamed"
	state := []tfstate.Resource{ecsStateResource(arn, "original-name")}
	live := []aws.ECSCluster{{ARN: arn, Name: "new-name", Status: "ACTIVE"}}

	results := DetectECSDrift(state, live)

	if len(results) != 1 || results[0].DriftType != "modified" {
		t.Errorf("expected 1 modified drift, got %+v", results)
	}
}

func TestDetectECSDrift_NoDrift(t *testing.T) {
	arn := "arn:aws:ecs:::cluster/stable"
	state := []tfstate.Resource{ecsStateResource(arn, "stable")}
	live := []aws.ECSCluster{{ARN: arn, Name: "stable", Status: "ACTIVE"}}

	results := DetectECSDrift(state, live)

	if len(results) != 0 {
		t.Errorf("expected no drift, got %+v", results)
	}
}
