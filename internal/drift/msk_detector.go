package drift

import (
	"fmt"

	"github.com/snyk/driftctl-lite/internal/aws"
	"github.com/snyk/driftctl-lite/internal/tfstate"
)

// DetectMSKDrift compares Terraform state MSK clusters against live AWS MSK clusters.
// It reports missing clusters, deleted/inactive clusters, and Kafka version mismatches.
func DetectMSKDrift(stateResources []tfstate.Resource, liveClusters []aws.MSKCluster) []DriftResult {
	var results []DriftResult

	liveByARN := make(map[string]aws.MSKCluster, len(liveClusters))
	for _, c := range liveClusters {
		liveByARN[c.ARN] = c
	}

	for _, res := range stateResources {
		if res.Type != "aws_msk_cluster" {
			continue
		}

		arn, _ := res.Attributes["arn"].(string)
		expectedVersion, _ := res.Attributes["kafka_version"].(string)

		live, found := liveByARN[arn]
		if !found {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   arn,
				DriftType:    "missing",
				Details:      "MSK cluster not found in AWS",
			})
			continue
		}

		if live.State == "DELETING" || live.State == "DELETED" {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   arn,
				DriftType:    "deleted",
				Details:      fmt.Sprintf("MSK cluster is in state %s", live.State),
			})
			continue
		}

		if expectedVersion != "" && live.KafkaVersion != expectedVersion {
			results = append(results, DriftResult{
				ResourceType: res.Type,
				ResourceID:   arn,
				DriftType:    "modified",
				Details:      fmt.Sprintf("kafka_version mismatch: state=%s live=%s", expectedVersion, live.KafkaVersion),
			})
		}
	}

	return results
}
