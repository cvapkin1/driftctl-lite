package drift

import (
	"testing"

	"github.com/driftctl-lite/internal/aws"
	"github.com/driftctl-lite/internal/tfstate"
	"github.com/stretchr/testify/assert"
)

func efsStateResource(id, throughputMode string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_efs_file_system",
		Name: id,
		Attributes: map[string]interface{}{
			"id":              id,
			"throughput_mode": throughputMode,
		},
	}
}

func TestDetectEFSDrift_Missing(t *testing.T) {
	state := []tfstate.Resource{efsStateResource("fs-missing", "bursting")}
	live := []aws.EFSFileSystem{}
	results := DetectEFSDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, "missing", results[0].Status)
	assert.Equal(t, "fs-missing", results[0].ResourceID)
}

func TestDetectEFSDrift_Deleted(t *testing.T) {
	state := []tfstate.Resource{efsStateResource("fs-del", "bursting")}
	live := []aws.EFSFileSystem{
		{ID: "fs-del", LifeCycleState: "deleting", ThroughputMode: "bursting"},
	}
	results := DetectEFSDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, "deleted", results[0].Status)
}

func TestDetectEFSDrift_ThroughputMismatch(t *testing.T) {
	state := []tfstate.Resource{efsStateResource("fs-tp", "bursting")}
	live := []aws.EFSFileSystem{
		{ID: "fs-tp", LifeCycleState: "available", ThroughputMode: "provisioned"},
	}
	results := DetectEFSDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, "drift", results[0].Status)
	assert.Contains(t, results[0].Message, "throughput_mode")
}

func TestDetectEFSDrift_NoDrift(t *testing.T) {
	state := []tfstate.Resource{efsStateResource("fs-ok", "bursting")}
	live := []aws.EFSFileSystem{
		{ID: "fs-ok", LifeCycleState: "available", ThroughputMode: "bursting"},
	}
	results := DetectEFSDrift(state, live)
	assert.Len(t, results, 1)
	assert.Equal(t, "ok", results[0].Status)
}
