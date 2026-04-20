package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

func eksStateResource(name, version string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_eks_cluster",
		Name: name,
		Attributes: map[string]interface{}{
			"name":    name,
			"version": version,
		},
	}
}

func TestDetectEKSDrift_Missing(t *testing.T) {
	resources := []tfstate.Resource{eksStateResource("prod-cluster", "1.28")}
	live := []aws.EKSCluster{}

	drifts := DetectEKSDrift(resources, live)
	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "not found in AWS")
}

func TestDetectEKSDrift_VersionMismatch(t *testing.T) {
	resources := []tfstate.Resource{eksStateResource("prod-cluster", "1.28")}
	live := []aws.EKSCluster{
		{Name: "prod-cluster", ARN: "arn:aws:eks:us-east-1:123:cluster/prod-cluster", Status: "ACTIVE", Version: "1.29"},
	}

	drifts := DetectEKSDrift(resources, live)
	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "version mismatch")
}

func TestDetectEKSDrift_DeletedStatus(t *testing.T) {
	resources := []tfstate.Resource{eksStateResource("prod-cluster", "1.28")}
	live := []aws.EKSCluster{
		{Name: "prod-cluster", ARN: "arn:aws:eks:us-east-1:123:cluster/prod-cluster", Status: "DELETING", Version: "1.28"},
	}

	drifts := DetectEKSDrift(resources, live)
	assert.Len(t, drifts, 1)
	assert.Contains(t, drifts[0], "unexpected status")
}

func TestDetectEKSDrift_NoDrift(t *testing.T) {
	resources := []tfstate.Resource{eksStateResource("prod-cluster", "1.29")}
	live := []aws.EKSCluster{
		{Name: "prod-cluster", ARN: "arn:aws:eks:us-east-1:123:cluster/prod-cluster", Status: "ACTIVE", Version: "1.29"},
	}

	drifts := DetectEKSDrift(resources, live)
	assert.Empty(t, drifts)
}
