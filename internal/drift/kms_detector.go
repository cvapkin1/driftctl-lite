package drift

import (
	"fmt"

	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

type KMSDriftResult struct {
	KeyID   string
	Issue   string
}

func DetectKMSDrift(resources []tfstate.Resource, liveKeys []aws.KMSKey) []KMSDriftResult {
	var results []KMSDriftResult

	liveByARN := make(map[string]aws.KMSKey)
	for _, k := range liveKeys {
		liveByARN[k.ARN] = k
	}

	for _, res := range resources {
		if res.Type != "aws_kms_key" {
			continue
		}

		arn, _ := res.Attributes["arn"].(string)
		keyID, _ := res.Attributes["key_id"].(string)

		live, found := liveByARN[arn]
		if !found {
			results = append(results, KMSDriftResult{
				KeyID: keyID,
				Issue: "KMS key not found in AWS (may have been deleted)",
			})
			continue
		}

		if live.State == "PendingDeletion" || live.State == "Disabled" {
			results = append(results, KMSDriftResult{
				KeyID: keyID,
				Issue: fmt.Sprintf("KMS key state is %q, expected Enabled", live.State),
			})
		}

		stateDesc, _ := res.Attributes["description"].(string)
		if stateDesc != "" && live.Description != stateDesc {
			results = append(results, KMSDriftResult{
				KeyID: keyID,
				Issue: fmt.Sprintf("description mismatch: state=%q live=%q", stateDesc, live.Description),
			})
		}
	}

	return results
}
