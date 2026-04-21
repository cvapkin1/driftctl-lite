package drift

import (
	"fmt"

	"github.com/owner/driftctl-lite/internal/aws"
)

type VPCDriftResult struct {
	ResourceID string
	DriftType  string
	Detail     string
}

// DetectVPCDrift compares Terraform state resources against live AWS VPCs.
// It reports missing VPCs and state mismatches.
func DetectVPCDrift(stateResources []map[string]interface{}, liveVPCs []aws.VPCResource) []VPCDriftResult {
	liveMap := make(map[string]aws.VPCResource, len(liveVPCs))
	for _, v := range liveVPCs {
		liveMap[v.ID] = v
	}

	var results []VPCDriftResult

	for _, res := range stateResources {
		id, _ := res["id"].(string)
		if id == "" {
			continue
		}

		live, found := liveMap[id]
		if !found {
			results = append(results, VPCDriftResult{
				ResourceID: id,
				DriftType:  "missing",
				Detail:     "VPC not found in AWS",
			})
			continue
		}

		if live.State != "available" {
			results = append(results, VPCDriftResult{
				ResourceID: id,
				DriftType:  "state_mismatch",
				Detail:     fmt.Sprintf("VPC state is %q, expected available", live.State),
			})
		}

		if cidr, ok := res["cidr_block"].(string); ok && cidr != "" && cidr != live.CIDR {
			results = append(results, VPCDriftResult{
				ResourceID: id,
				DriftType:  "cidr_mismatch",
				Detail:     fmt.Sprintf("CIDR in state %q differs from live %q", cidr, live.CIDR),
			})
		}
	}

	return results
}
