package drift

import (
	"testing"

	"github.com/acme/driftctl-lite/internal/aws"
	"github.com/acme/driftctl-lite/internal/tfstate"
	"github.com/stretchr/testify/assert"
)

func eipStateResource(allocationID, publicIP string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_eip",
		Name: allocationID,
		Attributes: map[string]interface{}{
			"allocation_id": allocationID,
			"public_ip":     publicIP,
		},
	}
}

func TestDetectElasticIPDrift_Missing(t *testing.T) {
	state := []tfstate.Resource{eipStateResource("eipalloc-abc", "1.2.3.4")}
	live := []aws.ElasticIPResource{}

	results := DetectElasticIPDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, DriftTypeMissing, results[0].DriftType)
	assert.Equal(t, "eipalloc-abc", results[0].ResourceID)
}

func TestDetectElasticIPDrift_PublicIPMismatch(t *testing.T) {
	state := []tfstate.Resource{eipStateResource("eipalloc-abc", "1.2.3.4")}
	live := []aws.ElasticIPResource{
		{AllocationID: "eipalloc-abc", PublicIP: "9.9.9.9", Domain: "vpc"},
	}

	results := DetectElasticIPDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, DriftTypeModified, results[0].DriftType)
	assert.Contains(t, results[0].Details, "public_ip mismatch")
}

func TestDetectElasticIPDrift_NoDrift(t *testing.T) {
	state := []tfstate.Resource{eipStateResource("eipalloc-abc", "1.2.3.4")}
	live := []aws.ElasticIPResource{
		{AllocationID: "eipalloc-abc", PublicIP: "1.2.3.4", Domain: "vpc"},
	}

	results := DetectElasticIPDrift(state, live)
	assert.Empty(t, results)
}

func TestDetectElasticIPDrift_IgnoresNonEIPResources(t *testing.T) {
	state := []tfstate.Resource{
		{
			Type: "aws_instance",
			Name: "web",
			Attributes: map[string]interface{}{"id": "i-abc"},
		},
	}
	live := []aws.ElasticIPResource{}

	results := DetectElasticIPDrift(state, live)
	assert.Empty(t, results)
}
