package drift

import (
	"fmt"

	"github.com/evilmartians/driftctl-lite/internal/aws"
)

// SecretsManagerDriftResult holds the result of a drift check for a single secret.
type SecretsManagerDriftResult struct {
	ResourceID string
	Status     string // "OK", "MISSING", "NAME_MISMATCH"
	Details    string
}

// DetectSecretsManagerDrift compares Terraform state resources of type
// "aws_secretsmanager_secret" against live secrets fetched from AWS.
// It returns a slice of drift results, one per state resource.
func DetectSecretsManagerDrift(stateResources []map[string]interface{}, liveSecrets []aws.SecretsManagerSecret) []SecretsManagerDriftResult {
	// Build a lookup map from secret ARN -> live secret for O(1) access.
	liveByARN := make(map[string]aws.SecretsManagerSecret, len(liveSecrets))
	for _, s := range liveSecrets {
		liveByARN[s.ARN] = s
	}

	var results []SecretsManagerDriftResult

	for _, res := range stateResources {
		resType, _ := res["type"].(string)
		if resType != "aws_secretsmanager_secret" {
			continue
		}

		attrs, _ := res["attributes"].(map[string]interface{})
		if attrs == nil {
			continue
		}

		arn, _ := attrs["arn"].(string)
		stateName, _ := attrs["name"].(string)

		if arn == "" {
			results = append(results, SecretsManagerDriftResult{
				ResourceID: stateName,
				Status:     "MISSING",
				Details:    "state resource has no ARN; cannot match against live secrets",
			})
			continue
		}

		live, found := liveByARN[arn]
		if !found {
			results = append(results, SecretsManagerDriftResult{
				ResourceID: arn,
				Status:     "MISSING",
				Details:    fmt.Sprintf("secret %q exists in state but was not found in AWS", arn),
			})
			continue
		}

		// Check that the secret name has not changed.
		if stateName != "" && live.Name != stateName {
			results = append(results, SecretsManagerDriftResult{
				ResourceID: arn,
				Status:     "NAME_MISMATCH",
				Details:    fmt.Sprintf("state name %q does not match live name %q", stateName, live.Name),
			})
			continue
		}

		results = append(results, SecretsManagerDriftResult{
			ResourceID: arn,
			Status:     "OK",
			Details:    "no drift detected",
		})
	}

	return results
}
