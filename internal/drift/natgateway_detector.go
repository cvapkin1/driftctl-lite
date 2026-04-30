package drift

import (
	"fmt"

	"github.com/jonhadfield/driftctl-lite/internal/aws"
)

type NATGatewayDriftResult struct {
	ResourceID string
	DriftType  string
	Details    string
}

func DetectNATGatewayDrift(stateResources []Resource, live []aws.NATGateway) []NATGatewayDriftResult {
	var results []NATGatewayDriftResult

	liveMap := make(map[string]aws.NATGateway)
	for _, gw := range live {
		liveMap[gw.ID] = gw
	}

	for _, res := range stateResources {
		if res.Type != "aws_nat_gateway" {
			continue
		}

		id, _ := res.Attributes["id"].(string)
		gw, found := liveMap[id]
		if !found {
			results = append(results, NATGatewayDriftResult{
				ResourceID: id,
				DriftType:  "missing",
				Details:    "NAT gateway not found in AWS",
			})
			continue
		}

		if gw.State == "deleted" || gw.State == "deleting" {
			results = append(results, NATGatewayDriftResult{
				ResourceID: id,
				DriftType:  "deleted",
				Details:    fmt.Sprintf("NAT gateway state is %s", gw.State),
			})
			continue
		}

		if stateSubnet, ok := res.Attributes["subnet_id"].(string); ok && stateSubnet != gw.SubnetID {
			results = append(results, NATGatewayDriftResult{
				ResourceID: id,
				DriftType:  "subnet_mismatch",
				Details:    fmt.Sprintf("expected subnet %s, got %s", stateSubnet, gw.SubnetID),
			})
		}
	}

	return results
}
