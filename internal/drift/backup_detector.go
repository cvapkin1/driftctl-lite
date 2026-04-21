package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

type BackupDriftResult struct {
	ResourceID string
	DriftType  string
	Details    string
}

func DetectBackupDrift(stateResources []tfstate.Resource, liveVaults []aws.BackupVault) []BackupDriftResult {
	var results []BackupDriftResult

	liveMap := make(map[string]aws.BackupVault)
	for _, v := range liveVaults {
		liveMap[v.Name] = v
	}

	for _, res := range stateResources {
		if res.Type != "aws_backup_vault" {
			continue
		}

		name, _ := res.Attributes["name"].(string)
		expectedARN, _ := res.Attributes["arn"].(string)

		live, found := liveMap[name]
		if !found {
			results = append(results, BackupDriftResult{
				ResourceID: name,
				DriftType:  "missing",
				Details:    fmt.Sprintf("backup vault %q not found in AWS", name),
			})
			continue
		}

		if expectedARN != "" && live.ARN != expectedARN {
			results = append(results, BackupDriftResult{
				ResourceID: name,
				DriftType:  "arn_mismatch",
				Details:    fmt.Sprintf("backup vault %q ARN mismatch: state=%q live=%q", name, expectedARN, live.ARN),
			})
		}
	}

	return results
}
