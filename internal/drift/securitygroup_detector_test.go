package drift

import (
	"testing"

	"github.com/owner/driftctl-lite/internal/aws"
	"github.com/owner/driftctl-lite/internal/tfstate"
	"github.com/stretchr/testify/assert"
)

func sgStateResource(id, description, vpcID string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_security_group",
		Attributes: map[string]interface{}{
			"id":          id,
			"description": description,
			"vpc_id":      vpcID,
		},
	}
}

func TestDetectSecurityGroupDrift_Missing(t *testing.T) {
	resources := []tfstate.Resource{sgStateResource("sg-abc", "web tier", "vpc-1")}
	results := DetectSecurityGroupDrift(resources, []aws.SecurityGroup{})

	assert.Len(t, results, 1)
	assert.Equal(t, "missing", results[0].DriftType)
	assert.Equal(t, "sg-abc", results[0].ResourceID)
}

func TestDetectSecurityGroupDrift_DescriptionMismatch(t *testing.T) {
	resources := []tfstate.Resource{sgStateResource("sg-abc", "web tier", "vpc-1")}
	live := []aws.SecurityGroup{
		{ID: "sg-abc", Description: "changed description", VPCID: "vpc-1"},
	}
	results := DetectSecurityGroupDrift(resources, live)

	assert.Len(t, results, 1)
	assert.Equal(t, "description_mismatch", results[0].DriftType)
}

func TestDetectSecurityGroupDrift_VPCMismatch(t *testing.T) {
	resources := []tfstate.Resource{sgStateResource("sg-abc", "web tier", "vpc-1")}
	live := []aws.SecurityGroup{
		{ID: "sg-abc", Description: "web tier", VPCID: "vpc-999"},
	}
	results := DetectSecurityGroupDrift(resources, live)

	assert.Len(t, results, 1)
	assert.Equal(t, "vpc_mismatch", results[0].DriftType)
}

func TestDetectSecurityGroupDrift_NoDrift(t *testing.T) {
	resources := []tfstate.Resource{sgStateResource("sg-abc", "web tier", "vpc-1")}
	live := []aws.SecurityGroup{
		{ID: "sg-abc", Description: "web tier", VPCID: "vpc-1"},
	}
	results := DetectSecurityGroupDrift(resources, live)

	assert.Empty(t, results)
}

func TestDetectSecurityGroupDrift_IgnoresOtherTypes(t *testing.T) {
	resources := []tfstate.Resource{
		{Type: "aws_instance", Attributes: map[string]interface{}{"id": "i-123"}},
	}
	results := DetectSecurityGroupDrift(resources, []aws.SecurityGroup{})

	assert.Empty(t, results)
}
