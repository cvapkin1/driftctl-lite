package drift

import (
	"testing"

	"github.com/driftctl-lite/internal/aws"
	"github.com/driftctl-lite/internal/tfstate"
	"github.com/stretchr/testify/assert"
)

func igwStateResource(id, vpcID string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_internet_gateway",
		Name: id,
		Attributes: map[string]interface{}{
			"id":     id,
			"vpc_id": vpcID,
		},
	}
}

func TestDetectInternetGatewayDrift_Missing(t *testing.T) {
	resources := []tfstate.Resource{igwStateResource("igw-aaa", "vpc-111")}
	results := DetectInternetGatewayDrift(resources, []aws.InternetGateway{})
	assert.Len(t, results, 1)
	assert.Equal(t, StatusMissing, results[0].Status)
	assert.Equal(t, "igw-aaa", results[0].ResourceID)
}

func TestDetectInternetGatewayDrift_Detached(t *testing.T) {
	resources := []tfstate.Resource{igwStateResource("igw-bbb", "vpc-222")}
	live := []aws.InternetGateway{
		{ID: "igw-bbb", VPCID: "", State: "detached"},
	}
	results := DetectInternetGatewayDrift(resources, live)
	assert.Len(t, results, 1)
	assert.Equal(t, StatusDrifted, results[0].Status)
	assert.Contains(t, results[0].Message, "detached")
}

func TestDetectInternetGatewayDrift_VPCMismatch(t *testing.T) {
	resources := []tfstate.Resource{igwStateResource("igw-ccc", "vpc-333")}
	live := []aws.InternetGateway{
		{ID: "igw-ccc", VPCID: "vpc-999", State: "attached"},
	}
	results := DetectInternetGatewayDrift(resources, live)
	assert.Len(t, results, 1)
	assert.Equal(t, StatusDrifted, results[0].Status)
	assert.Contains(t, results[0].Message, "vpc_id mismatch")
}

func TestDetectInternetGatewayDrift_NoDrift(t *testing.T) {
	resources := []tfstate.Resource{igwStateResource("igw-ddd", "vpc-444")}
	live := []aws.InternetGateway{
		{ID: "igw-ddd", VPCID: "vpc-444", State: "attached"},
	}
	results := DetectInternetGatewayDrift(resources, live)
	assert.Len(t, results, 1)
	assert.Equal(t, StatusOK, results[0].Status)
}

func TestDetectInternetGatewayDrift_IgnoresOtherTypes(t *testing.T) {
	resources := []tfstate.Resource{
		{Type: "aws_vpc", Name: "main", Attributes: map[string]interface{}{"id": "vpc-111"}},
	}
	results := DetectInternetGatewayDrift(resources, []aws.InternetGateway{})
	assert.Empty(t, results)
}
