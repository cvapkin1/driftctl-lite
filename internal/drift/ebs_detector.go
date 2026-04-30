package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

const (
	ebsDeletedState = "deleted"
)

// DetectEBSDrift compares Terraform state EBS volumes against live AWS volumes.
func DetectEBSDrift(resources []tfstate.Resource, live []aws.EBSVolume) []DriftResult {
	liveMap := make(map[string]aws.EBSVolume, len(live))
	for _, v := range live {
		liveMap[v.ID] = v
	}

	var results []DriftResult

	for _, r := range resources {
		if r.Type != "aws_ebs_volume" {
			continue
		}

		volumeID, _ := r.Attributes["id"].(string)
		liveVol, found := liveMap[volumeID]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   volumeID,
				ResourceType: r.Type,
				DriftType:    DriftTypeMissing,
				Details:      fmt.Sprintf("EBS volume %s not found in AWS", volumeID),
			})
			continue
		}

		if liveVol.State == ebsDeletedState {
			results = append(results, DriftResult{
				ResourceID:   volumeID,
				ResourceType: r.Type,
				DriftType:    DriftTypeDeleted,
				Details:      fmt.Sprintf("EBS volume %s is in deleted state", volumeID),
			})
			continue
		}

		expectedType, _ := r.Attributes["volume_type"].(string)
		if expectedType != "" && liveVol.VolumeType != expectedType {
			results = append(results, DriftResult{
				ResourceID:   volumeID,
				ResourceType: r.Type,
				DriftType:    DriftTypeModified,
				Details:      fmt.Sprintf("EBS volume %s type mismatch: state=%s live=%s", volumeID, expectedType, liveVol.VolumeType),
			})
		}
	}

	return results
}
