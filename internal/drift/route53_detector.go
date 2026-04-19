package drift

import (
	"fmt"

	"github.com/owner/driftctl-lite/internal/aws"
)

type Route53DriftResult struct {
	ZoneID string
	Issue  string
}

func DetectRoute53Drift(stateZones []map[string]interface{}, liveZones []aws.Route53Zone) []Route53DriftResult {
	liveIndex := make(map[string]aws.Route53Zone, len(liveZones))
	for _, z := range liveZones {
		liveIndex[z.ID] = z
	}

	var results []Route53DriftResult

	for _, s := range stateZones {
		id, _ := s["zone_id"].(string)
		name, _ := s["name"].(string)

		live, found := liveIndex[id]
		if !found {
			results = append(results, Route53DriftResult{
				ZoneID: id,
				Issue:  "hosted zone not found in AWS",
			})
			continue
		}

		if name != "" && live.Name != name {
			results = append(results, Route53DriftResult{
				ZoneID: id,
				Issue:  fmt.Sprintf("name mismatch: state=%s live=%s", name, live.Name),
			})
		}
	}

	return results
}
