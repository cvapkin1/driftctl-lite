package drift

import (
	"fmt"

	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

// DetectElasticsearchDrift compares Terraform state ES domains against live AWS resources.
func DetectElasticsearchDrift(stateResources []tfstate.Resource, liveResources []aws.ESResource) []DriftResult {
	liveMap := make(map[string]aws.ESResource)
	for _, r := range liveResources {
		liveMap[r.Name] = r
	}

	var results []DriftResult

	for _, sr := range stateResources {
		if sr.Type != "aws_elasticsearch_domain" {
			continue
		}

		domainName, _ := sr.Attributes["domain_name"].(string)
		if domainName == "" {
			domainName, _ = sr.Attributes["id"].(string)
		}

		live, found := liveMap[domainName]
		if !found {
			results = append(results, DriftResult{
				ResourceType: sr.Type,
				ResourceID:   domainName,
				Status:       StatusMissing,
				Message:      fmt.Sprintf("Elasticsearch domain %q not found in AWS", domainName),
			})
			continue
		}

		if live.Deleted {
			results = append(results, DriftResult{
				ResourceType: sr.Type,
				ResourceID:   domainName,
				Status:       StatusDeleted,
				Message:      fmt.Sprintf("Elasticsearch domain %q is marked as deleted in AWS", domainName),
			})
			continue
		}

		stateVersion, _ := sr.Attributes["elasticsearch_version"].(string)
		if stateVersion != "" && stateVersion != live.Version {
			results = append(results, DriftResult{
				ResourceType: sr.Type,
				ResourceID:   domainName,
				Status:       StatusModified,
				Message:      fmt.Sprintf("Elasticsearch domain %q version mismatch: state=%s live=%s", domainName, stateVersion, live.Version),
			})
			continue
		}

		results = append(results, DriftResult{
			ResourceType: sr.Type,
			ResourceID:   domainName,
			Status:       StatusOK,
			Message:      fmt.Sprintf("Elasticsearch domain %q is in sync", domainName),
		})
	}

	return results
}
