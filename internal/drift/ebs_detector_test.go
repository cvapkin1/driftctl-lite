package drift

import (
	"testing"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
	"github.com/stretchr/testify/assert"
)

func ebsStateResource(id, volumeType string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_ebs_volume",
		Name: id,
		Attributes: map[string]interface{}{
			"id":          id,
			"volume_type": volumeType,
		},
	}
}

func TestDetectEBSDrift_Missing(t *testing.T) {
	resources := []tfstate.Resource{ebsStateResource("vol-aaa", "gp3")}
	results := DetectEBSDrift(resources, []aws.EBSVolume{})
	assert.Len(t, results, 1)
	assert.Equal(t, DriftTypeMissing, results[0].DriftType)
	assert.Equal(t, "vol-aaa", results[0].ResourceID)
}

func TestDetectEBSDrift_Deleted(t *testing.T) {
	resources := []tfstate.Resource{ebsStateResource("vol-bbb", "gp2")}
	live := []aws.EBSVolume{
		{ID: "vol-bbb", State: "deleted", VolumeType: "gp2"},
	}
	results := DetectEBSDrift(resources, live)
	assert.Len(t, results, 1)
	assert.Equal(t, DriftTypeDeleted, results[0].DriftType)
}

func TestDetectEBSDrift_VolumeTypeMismatch(t *testing.T) {
	resources := []tfstate.Resource{ebsStateResource("vol-ccc", "gp3")}
	live := []aws.EBSVolume{
		{ID: "vol-ccc", State: "available", VolumeType: "gp2"},
	}
	results := DetectEBSDrift(resources, live)
	assert.Len(t, results, 1)
	assert.Equal(t, DriftTypeModified, results[0].DriftType)
	assert.Contains(t, results[0].Details, "gp3")
	assert.Contains(t, results[0].Details, "gp2")
}

func TestDetectEBSDrift_NoDrift(t *testing.T) {
	resources := []tfstate.Resource{ebsStateResource("vol-ddd", "io1")}
	live := []aws.EBSVolume{
		{ID: "vol-ddd", State: "available", VolumeType: "io1"},
	}
	results := DetectEBSDrift(resources, live)
	assert.Empty(t, results)
}

func TestDetectEBSDrift_IgnoresNonEBSResources(t *testing.T) {
	r := tfstate.Resource{
		Type: "aws_instance",
		Name: "web",
		Attributes: map[string]interface{}{"id": "i-123"},
	}
	results := DetectEBSDrift([]tfstate.Resource{r}, []aws.EBSVolume{})
	assert.Empty(t, results)
}
