package drift

import (
	"fmt"

	"github.com/driftctl-lite/internal/aws"
	"github.com/driftctl-lite/internal/tfstate"
)

type EFSDriftResult struct {
	ResourceID string
	Status     string
	Message    string
}

func DetectEFSDrift(stateResources []tfstate.Resource, live []aws.EFSFileSystem) []EFSDriftResult {
	var results []EFSDriftResult

	liveMap := make(map[string]aws.EFSFileSystem)
	for _, fs := range live {
		liveMap[fs.ID] = fs
	}

	for _, res := range stateResources {
		if res.Type != "aws_efs_file_system" {
			continue
		}
		id, _ := res.Attributes["id"].(string)
		liveFS, found := liveMap[id]
		if !found {
			results = append(results, EFSDriftResult{
				ResourceID: id,
				Status:     "missing",
				Message:    fmt.Sprintf("EFS file system %s not found in AWS", id),
			})
			continue
		}
		if liveFS.LifeCycleState == "deleted" || liveFS.LifeCycleState == "deleting" {
			results = append(results, EFSDriftResult{
				ResourceID: id,
				Status:     "deleted",
				Message:    fmt.Sprintf("EFS file system %s is in state %s", id, liveFS.LifeCycleState),
			})
			continue
		}
		stateThroughput, _ := res.Attributes["throughput_mode"].(string)
		if stateThroughput != "" && stateThroughput != liveFS.ThroughputMode {
			results = append(results, EFSDriftResult{
				ResourceID: id,
				Status:     "drift",
				Message:    fmt.Sprintf("EFS file system %s throughput_mode changed: state=%s live=%s", id, stateThroughput, liveFS.ThroughputMode),
			})
			continue
		}
		results = append(results, EFSDriftResult{
			ResourceID: id,
			Status:     "ok",
			Message:    fmt.Sprintf("EFS file system %s is in sync", id),
		})
	}
	return results
}
