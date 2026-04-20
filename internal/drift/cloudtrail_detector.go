package drift

import (
	"fmt"

	"github.com/acme/driftctl-lite/internal/aws"
	"github.com/acme/driftctl-lite/internal/tfstate"
)

// DetectCloudTrailDrift compares CloudTrail trails in Terraform state against live AWS resources.
func DetectCloudTrailDrift(resources []tfstate.Resource, live []aws.CloudTrailResource) []DriftResult {
	var results []DriftResult

	liveIndex := make(map[string]aws.CloudTrailResource, len(live))
	for _, trail := range live {
		liveIndex[trail.Name] = trail
	}

	for _, res := range resources {
		if res.Type != "aws_cloudtrail" {
			continue
		}

		name, _ := res.Attributes["name"].(string)
		expectedBucket, _ := res.Attributes["s3_bucket_name"].(string)

		liveTrail, found := liveIndex[name]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   name,
				ResourceType: "aws_cloudtrail",
				DriftType:    "missing",
				Details:      fmt.Sprintf("trail %q not found in AWS", name),
			})
			continue
		}

		if expectedBucket != "" && liveTrail.S3Bucket != expectedBucket {
			results = append(results, DriftResult{
				ResourceID:   name,
				ResourceType: "aws_cloudtrail",
				DriftType:    "modified",
				Details:      fmt.Sprintf("s3_bucket_name mismatch: state=%q live=%q", expectedBucket, liveTrail.S3Bucket),
			})
		}
	}

	return results
}
