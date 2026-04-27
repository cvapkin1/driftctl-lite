package drift

import (
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

// DetectCloudFrontDrift compares CloudFront distributions from Terraform state
// against live AWS resources and returns a list of DriftResult entries.
func DetectCloudFrontDrift(state []tfstate.Resource, live []map[string]string) []DriftResult {
	liveMap := make(map[string]map[string]string, len(live))
	for _, dist := range live {
		if id, ok := dist["id"]; ok {
			liveMap[id] = dist
		}
	}

	var results []DriftResult

	for _, res := range state {
		if res.Type != "aws_cloudfront_distribution" {
			continue
		}

		id := res.Attributes["id"]
		liveDist, found := liveMap[id]
		if !found {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   id,
				Status:       "missing",
				Message:      "CloudFront distribution not found in AWS",
			})
			continue
		}

		wantDomain := res.Attributes["domain_name"]
		gotDomain := liveDist["domain_name"]
		if wantDomain != gotDomain {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   id,
				Status:       "changed",
				Message:      "domain_name mismatch: state=" + wantDomain + " live=" + gotDomain,
			})
			continue
		}

		wantStatus := res.Attributes["status"]
		gotStatus := liveDist["status"]
		if wantStatus != gotStatus {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   id,
				Status:       "changed",
				Message:      "status mismatch: state=" + wantStatus + " live=" + gotStatus,
			})
			continue
		}
	}

	return results
}
