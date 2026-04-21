package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

type OpenSearchDriftResult struct {
	ResourceID string
	DriftType  string
	Details    string
}

func DetectOpenSearchDrift(stateResources []tfstate.Resource, liveDomains []aws.OpenSearchDomain) []OpenSearchDriftResult {
	var results []OpenSearchDriftResult

	liveMap := make(map[string]aws.OpenSearchDomain, len(liveDomains))
	for _, d := range liveDomains {
		liveMap[d.DomainName] = d
	}

	for _, res := range stateResources {
		if res.Type != "aws_opensearch_domain" {
			continue
		}

		domainName, _ := res.Attributes["domain_name"].(string)
		if domainName == "" {
			domainName, _ = res.Attributes["id"].(string)
		}

		live, found := liveMap[domainName]
		if !found {
			results = append(results, OpenSearchDriftResult{
				ResourceID: domainName,
				DriftType:  "missing",
				Details:    "OpenSearch domain not found in AWS",
			})
			continue
		}

		if live.Deleted {
			results = append(results, OpenSearchDriftResult{
				ResourceID: domainName,
				DriftType:  "deleted",
				Details:    "OpenSearch domain is marked as deleted in AWS",
			})
			continue
		}

		stateVersion, _ := res.Attributes["engine_version"].(string)
		if stateVersion != "" && stateVersion != live.EngineVersion {
			results = append(results, OpenSearchDriftResult{
				ResourceID: domainName,
				DriftType:  "engine_version_mismatch",
				Details:    fmt.Sprintf("expected engine version %q, got %q", stateVersion, live.EngineVersion),
			})
		}
	}

	return results
}
