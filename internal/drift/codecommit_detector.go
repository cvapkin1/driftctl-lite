package drift

import (
	"fmt"

	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

// DetectCodeCommitDrift compares Terraform state repositories against live AWS CodeCommit repositories.
// It reports repositories that are missing or have a mismatched ARN.
func DetectCodeCommitDrift(stateResources []tfstate.Resource, liveRepos []aws.CodeCommitRepository) []DriftResult {
	var results []DriftResult

	liveIndex := make(map[string]aws.CodeCommitRepository, len(liveRepos))
	for _, r := range liveRepos {
		liveIndex[r.Name] = r
	}

	for _, res := range stateResources {
		if res.Type != "aws_codecommit_repository" {
			continue
		}

		stateName, _ := res.Attributes["repository_name"].(string)
		stateARN, _ := res.Attributes["arn"].(string)

		live, found := liveIndex[stateName]
		if !found {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   stateName,
				DriftType:    DriftTypeMissing,
				Details:      fmt.Sprintf("repository %q not found in AWS", stateName),
			})
			continue
		}

		if stateARN != "" && live.ARN != stateARN {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   stateName,
				DriftType:    DriftTypeModified,
				Details:      fmt.Sprintf("ARN mismatch: state=%q live=%q", stateARN, live.ARN),
			})
		}
	}

	return results
}
