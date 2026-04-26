package drift

import (
	"fmt"

	"github.com/acme/driftctl-lite/internal/aws"
)

type ElasticBeanstalkDriftResult struct {
	EnvironmentID   string
	DriftType       string
	Details         string
}

// DetectElasticBeanstalkDrift compares Terraform state resources against live
// Elastic Beanstalk environments and returns any detected drift.
func DetectElasticBeanstalkDrift(
	stateResources []map[string]interface{},
	live []aws.ElasticBeanstalkEnvironment,
) []ElasticBeanstalkDriftResult {
	liveByID := make(map[string]aws.ElasticBeanstalkEnvironment, len(live))
	for _, e := range live {
		liveByID[e.ID] = e
	}

	var results []ElasticBeanstalkDriftResult

	for _, res := range stateResources {
		envID, _ := res["id"].(string)
		if envID == "" {
			continue
		}

		liveEnv, found := liveByID[envID]
		if !found {
			results = append(results, ElasticBeanstalkDriftResult{
				EnvironmentID: envID,
				DriftType:     "missing",
				Details:       "environment not found in AWS",
			})
			continue
		}

		if liveEnv.Status == "Terminated" || liveEnv.Status == "Terminating" {
			results = append(results, ElasticBeanstalkDriftResult{
				EnvironmentID: envID,
				DriftType:     "terminated",
				Details:       fmt.Sprintf("environment status is %s", liveEnv.Status),
			})
			continue
		}

		if stateCNAME, ok := res["cname"].(string); ok && stateCNAME != "" {
			if stateCNAME != liveEnv.CNAME {
				results = append(results, ElasticBeanstalkDriftResult{
					EnvironmentID: envID,
					DriftType:     "cname_mismatch",
					Details:       fmt.Sprintf("state CNAME %q differs from live %q", stateCNAME, liveEnv.CNAME),
				})
			}
		}
	}

	return results
}
