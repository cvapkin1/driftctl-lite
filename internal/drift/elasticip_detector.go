package drift

import (
	"fmt"

	"github.com/acme/driftctl-lite/internal/aws"
	"github.com/acme/driftctl-lite/internal/tfstate"
)

// DetectElasticIPDrift compares Elastic IP resources from Terraform state
// against live AWS resources and returns a list of drift findings.
func DetectElasticIPDrift(stateResources []tfstate.Resource, liveAddresses []aws.ElasticIPResource) []DriftResult {
	liveMap := make(map[string]aws.ElasticIPResource, len(liveAddresses))
	for _, addr := range liveAddresses {
		liveMap[addr.AllocationID] = addr
	}

	var results []DriftResult

	for _, res := range stateResources {
		if res.Type != "aws_eip" {
			continue
		}

		allocationID, _ := res.Attributes["allocation_id"].(string)
		if allocationID == "" {
			allocationID, _ = res.Attributes["id"].(string)
		}

		live, found := liveMap[allocationID]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   allocationID,
				ResourceType: res.Type,
				DriftType:    DriftTypeMissing,
				Details:      fmt.Sprintf("Elastic IP %s not found in AWS", allocationID),
			})
			continue
		}

		expectedPublicIP, _ := res.Attributes["public_ip"].(string)
		if expectedPublicIP != "" && live.PublicIP != expectedPublicIP {
			results = append(results, DriftResult{
				ResourceID:   allocationID,
				ResourceType: res.Type,
				DriftType:    DriftTypeModified,
				Details:      fmt.Sprintf("public_ip mismatch: state=%s live=%s", expectedPublicIP, live.PublicIP),
			})
			continue
		}
	}

	return results
}
