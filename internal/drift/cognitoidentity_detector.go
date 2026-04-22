package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

type CognitoIdentityDrift struct {
	PoolID string
	Reason string
}

func DetectCognitoIdentityDrift(resources []tfstate.Resource, pools []aws.CognitoIdentityPool) []CognitoIdentityDrift {
	var drifts []CognitoIdentityDrift

	liveMap := make(map[string]aws.CognitoIdentityPool)
	for _, p := range pools {
		liveMap[p.ID] = p
	}

	for _, res := range resources {
		if res.Type != "aws_cognito_identity_pool" {
			continue
		}

		poolID, _ := res.Attributes["id"].(string)
		expectedName, _ := res.Attributes["identity_pool_name"].(string)

		live, found := liveMap[poolID]
		if !found {
			drifts = append(drifts, CognitoIdentityDrift{
				PoolID: poolID,
				Reason: "identity pool not found in AWS",
			})
			continue
		}

		if expectedName != "" && live.Name != expectedName {
			drifts = append(drifts, CognitoIdentityDrift{
				PoolID: poolID,
				Reason: fmt.Sprintf("name mismatch: state=%q live=%q", expectedName, live.Name),
			})
		}
	}

	return drifts
}
