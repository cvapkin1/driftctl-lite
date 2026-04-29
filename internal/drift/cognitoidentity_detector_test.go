package drift_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"driftctl-lite/internal/drift"
	"driftctl-lite/internal/tfstate"
)

func cognitoStateResource(id, name string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_cognito_identity_pool",
		Name: name,
		Attributes: map[string]interface{}{
			"id":                 id,
			"identity_pool_name": name,
		},
	}
}

func TestDetectCognitoIdentityDrift_Missing(t *testing.T) {
	state := []tfstate.Resource{
		cognitoStateResource("us-east-1:abc-123", "my-pool"),
	}
	live := []map[string]string{}

	results := drift.DetectCognitoIdentityDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, "missing", results[0].Status)
	assert.Equal(t, "us-east-1:abc-123", results[0].ResourceID)
}

func TestDetectCognitoIdentityDrift_NameMismatch(t *testing.T) {
	state := []tfstate.Resource{
		cognitoStateResource("us-east-1:abc-123", "my-pool"),
	}
	live := []map[string]string{
		{"id": "us-east-1:abc-123", "name": "renamed-pool"},
	}

	results := drift.DetectCognitoIdentityDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, "drift", results[0].Status)
	assert.Contains(t, results[0].Detail, "name mismatch")
}

func TestDetectCognitoIdentityDrift_NoDrift(t *testing.T) {
	state := []tfstate.Resource{
		cognitoStateResource("us-east-1:abc-123", "my-pool"),
	}
	live := []map[string]string{
		{"id": "us-east-1:abc-123", "name": "my-pool"},
	}

	results := drift.DetectCognitoIdentityDrift(state, live)
	assert.Empty(t, results)
}
