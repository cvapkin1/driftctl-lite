package drift

import (
	"fmt"

	"github.com/driftctl-lite/internal/aws"
	"github.com/driftctl-lite/internal/tfstate"
)

// DetectInternetGatewayDrift compares Terraform state internet gateways against live AWS resources.
func DetectInternetGatewayDrift(resources []tfstate.Resource, live []aws.InternetGateway) []DriftResult {
	liveMap := make(map[string]aws.InternetGateway)
	for _, igw := range live {
		liveMap[igw.ID] = igw
	}

	var results []DriftResult

	for _, res := range resources {
		if res.Type != "aws_internet_gateway" {
			continue
		}

		igwID, _ := res.Attributes["id"].(string)
		liveIGW, found := liveMap[igwID]

		if !found {
			results = append(results, DriftResult{
				ResourceID:   igwID,
				ResourceType: res.Type,
				Status:       StatusMissing,
				Message:      fmt.Sprintf("internet gateway %s not found in AWS", igwID),
			})
			continue
		}

		if liveIGW.State == "detached" {
			results = append(results, DriftResult{
				ResourceID:   igwID,
				ResourceType: res.Type,
				Status:       StatusDrifted,
				Message:      fmt.Sprintf("internet gateway %s is detached", igwID),
			})
			continue
		}

		stateVPCID, _ := res.Attributes["vpc_id"].(string)
		if stateVPCID != "" && liveIGW.VPCID != stateVPCID {
			results = append(results, DriftResult{
				ResourceID:   igwID,
				ResourceType: res.Type,
				Status:       StatusDrifted,
				Message:      fmt.Sprintf("internet gateway %s vpc_id mismatch: state=%s live=%s", igwID, stateVPCID, liveIGW.VPCID),
			})
			continue
		}

		results = append(results, DriftResult{
			ResourceID:   igwID,
			ResourceType: res.Type,
			Status:       StatusOK,
			Message:      fmt.Sprintf("internet gateway %s is in sync", igwID),
		})
	}

	return results
}
