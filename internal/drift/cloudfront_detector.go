package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

// DetectCloudFrontDrift compares Terraform state CloudFront distributions
// against live AWS distributions and returns a list of drift findings.
func DetectCloudFrontDrift(stateResources []tfstate.Resource, live []aws.CloudFrontResource) []string {
	var drifts []string

	liveMap := make(map[string]aws.CloudFrontResource, len(live))
	for _, d := range live {
		liveMap[d.ID] = d
	}

	for _, sr := range stateResources {
		if sr.Type != "aws_cloudfront_distribution" {
			continue
		}

		id, _ := sr.Attributes["id"].(string)
		if id == "" {
			continue
		}

		liveD, found := liveMap[id]
		if !found {
			drifts = append(drifts, fmt.Sprintf("CloudFront distribution %s missing in AWS", id))
			continue
		}

		if liveD.Status == "Deployed" == false && liveD.Status != "" {
			if liveD.Status == "InProgress" {
				// still deploying — not a drift
				continue
			}
		}

		stateDomain, _ := sr.Attributes["domain_name"].(string)
		if stateDomain != "" && stateDomain != liveD.DomainName {
			drifts = append(drifts, fmt.Sprintf(
				"CloudFront distribution %s domain mismatch: state=%s live=%s",
				id, stateDomain, liveD.DomainName,
			))
		}

		stateEnabled, hasEnabled := sr.Attributes["enabled"].(bool)
		if hasEnabled && stateEnabled != liveD.Enabled {
			drifts = append(drifts, fmt.Sprintf(
				"CloudFront distribution %s enabled mismatch: state=%v live=%v",
				id, stateEnabled, liveD.Enabled,
			))
		}
	}

	return drifts
}
