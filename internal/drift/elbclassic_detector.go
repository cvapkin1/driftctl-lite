package drift

import (
	"fmt"

	internalaws "driftctl-lite/internal/aws"
	"driftctl-lite/internal/tfstate"
)

// DetectELBClassicDrift compares Terraform state classic ELBs against live AWS resources.
func DetectELBClassicDrift(resources []tfstate.Resource, live []internalaws.ELBClassicResource) []DriftResult {
	liveMap := make(map[string]internalaws.ELBClassicResource)
	for _, lb := range live {
		liveMap[lb.Name] = lb
	}

	var results []DriftResult

	for _, res := range resources {
		if res.Type != "aws_elb" {
			continue
		}

		name, _ := res.Attributes["name"].(string)
		lb, found := liveMap[name]

		if !found {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   name,
				Status:       StatusMissing,
				Details:      fmt.Sprintf("ELB classic load balancer %q not found in AWS", name),
			})
			continue
		}

		if lb.State != "active" {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   name,
				Status:       StatusModified,
				Details:      fmt.Sprintf("ELB classic load balancer %q has unexpected state: %q", name, lb.State),
			})
			continue
		}

		expectedDNS, _ := res.Attributes["dns_name"].(string)
		if expectedDNS != "" && lb.DNSName != expectedDNS {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   name,
				Status:       StatusModified,
				Details:      fmt.Sprintf("ELB classic load balancer %q DNS mismatch: state=%q live=%q", name, expectedDNS, lb.DNSName),
			})
			continue
		}

		results = append(results, DriftResult{
			ResourceType: res.Type,
			ResourceID:   name,
			Status:       StatusOK,
			Details:      fmt.Sprintf("ELB classic load balancer %q is in sync", name),
		})
	}

	return results
}
