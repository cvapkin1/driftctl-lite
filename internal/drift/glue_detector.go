package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

// DetectGlueDrift compares Glue jobs in Terraform state against live AWS resources.
// It reports jobs that are missing or have a role mismatch.
func DetectGlueDrift(resources []tfstate.Resource, liveJobs []aws.GlueJob) []string {
	var drifts []string

	liveMap := make(map[string]aws.GlueJob, len(liveJobs))
	for _, j := range liveJobs {
		liveMap[j.Name] = j
	}

	for _, r := range resources {
		if r.Type != "aws_glue_job" {
			continue
		}

		name, _ := r.Attributes["name"].(string)
		if name == "" {
			continue
		}

		live, found := liveMap[name]
		if !found {
			drifts = append(drifts, fmt.Sprintf("[MISSING] aws_glue_job %q not found in AWS", name))
			continue
		}

		if stateRole, ok := r.Attributes["role_arn"].(string); ok && stateRole != "" {
			if live.Role != stateRole {
				drifts = append(drifts, fmt.Sprintf(
					"[CHANGED] aws_glue_job %q role mismatch: state=%q live=%q",
					name, stateRole, live.Role,
				))
			}
		}

		if stateVersion, ok := r.Attributes["glue_version"].(string); ok && stateVersion != "" {
			if live.GlueVersion != stateVersion {
				drifts = append(drifts, fmt.Sprintf(
					"[CHANGED] aws_glue_job %q glue_version mismatch: state=%q live=%q",
					name, stateVersion, live.GlueVersion,
				))
			}
		}
	}

	return drifts
}
