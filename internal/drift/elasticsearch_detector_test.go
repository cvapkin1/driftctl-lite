package drift

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

func esStateResource(domainName, version string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_elasticsearch_domain",
		Name: domainName,
		Attributes: map[string]interface{}{
			"domain_name":           domainName,
			"elasticsearch_version": version,
		},
	}
}

func TestDetectElasticsearchDrift_Missing(t *testing.T) {
	state := []tfstate.Resource{esStateResource("prod-es", "7.10")}
	live := []aws.ESResource{}

	results := DetectElasticsearchDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, StatusMissing, results[0].Status)
	assert.Equal(t, "prod-es", results[0].ResourceID)
}

func TestDetectElasticsearchDrift_Deleted(t *testing.T) {
	state := []tfstate.Resource{esStateResource("prod-es", "7.10")}
	live := []aws.ESResource{
		{Name: "prod-es", Version: "7.10", Deleted: true},
	}

	results := DetectElasticsearchDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, StatusDeleted, results[0].Status)
}

func TestDetectElasticsearchDrift_VersionMismatch(t *testing.T) {
	state := []tfstate.Resource{esStateResource("prod-es", "7.10")}
	live := []aws.ESResource{
		{Name: "prod-es", Version: "6.8", Deleted: false},
	}

	results := DetectElasticsearchDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, StatusModified, results[0].Status)
	assert.Contains(t, results[0].Message, "version mismatch")
}

func TestDetectElasticsearchDrift_NoDrift(t *testing.T) {
	state := []tfstate.Resource{esStateResource("prod-es", "7.10")}
	live := []aws.ESResource{
		{Name: "prod-es", Version: "7.10", Deleted: false},
	}

	results := DetectElasticsearchDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, StatusOK, results[0].Status)
}
