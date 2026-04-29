package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

// DetectSubnetDrift compares Terraform state subnet resources against live AWS subnets.
func DetectSubnetDrift(resources []tfstate.Resource, live []aws.SubnetResource) []DriftResult {
	liveMap := make(map[string]aws.SubnetResource, len(live))
	for _, s := range live {
		liveMap[s.ID] = s
	}

	var results []DriftResult

	for _, res := range resources {
		if res.Type != "aws_subnet" {
			continue
		}

		subnetID, _ := res.Attributes["id"].(string)
		liveSubnet, found := liveMap[subnetID]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   subnetID,
				ResourceType: res.Type,
				DriftType:    "missing",
				Details:      "subnet not found in AWS",
			})
			continue
		}

		if liveSubnet.State != "available" {
			results = append(results, DriftResult{
				ResourceID:   subnetID,
				ResourceType: res.Type,
				DriftType:    "deleted",
				Details:      fmt.Sprintf("subnet state is %q", liveSubnet.State),
			})
			continue
		}

		stateCIDR, _ := res.Attributes["cidr_block"].(string)
		if stateCIDR != "" && stateCIDR != liveSubnet.CIDRBlock {
			results = append(results, DriftResult{
				ResourceID:   subnetID,
				ResourceType: res.Type,
				DriftType:    "modified",
				Details:      fmt.Sprintf("cidr_block mismatch: state=%q live=%q", stateCIDR, liveSubnet.CIDRBlock),
			})
			continue
		}

		stateVPC, _ := res.Attributes["vpc_id"].(string)
		if stateVPC != "" && stateVPC != liveSubnet.VPCID {
			results = append(results, DriftResult{
				ResourceID:   subnetID,
				ResourceType: res.Type,
				DriftType:    "modified",
				Details:      fmt.Sprintf("vpc_id mismatch: state=%q live=%q", stateVPC, liveSubnet.VPCID),
			})
		}
	}

	return results
}
