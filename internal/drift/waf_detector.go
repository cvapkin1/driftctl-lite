package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

// DetectWAFDrift compares WAF Web ACLs from Terraform state against live AWS resources.
// It reports resources that are missing or have a mismatched ARN.
func DetectWAFDrift(resources []tfstate.Resource, liveACLs []aws.WAFWebACL) []string {
	var drifts []string

	liveByID := make(map[string]aws.WAFWebACL, len(liveACLs))
	for _, acl := range liveACLs {
		liveByID[acl.ID] = acl
	}

	for _, res := range resources {
		if res.Type != "aws_wafv2_web_acl" {
			continue
		}

		id, _ := res.Attributes["id"].(string)
		if id == "" {
			drifts = append(drifts, fmt.Sprintf("[WAF] resource %q has no id in state", res.Name))
			continue
		}

		live, found := liveByID[id]
		if !found {
			drifts = append(drifts, fmt.Sprintf("[WAF] web ACL %q (id=%s) not found in AWS", res.Name, id))
			continue
		}

		stateARN, _ := res.Attributes["arn"].(string)
		if stateARN != "" && stateARN != live.ARN {
			drifts = append(drifts, fmt.Sprintf(
				"[WAF] web ACL %q ARN mismatch: state=%s live=%s",
				res.Name, stateARN, live.ARN,
			))
		}

		stateName, _ := res.Attributes["name"].(string)
		if stateName != "" && stateName != live.Name {
			drifts = append(drifts, fmt.Sprintf(
				"[WAF] web ACL %q name mismatch: state=%s live=%s",
				res.Name, stateName, live.Name,
			))
		}
	}

	return drifts
}
