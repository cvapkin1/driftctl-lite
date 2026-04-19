package drift

import (
	"fmt"

	"github.com/owner/driftctl-lite/internal/aws"
	"github.com/owner/driftctl-lite/internal/tfstate"
)

// DetectS3Drift compares S3 buckets in Terraform state against live AWS buckets.
func DetectS3Drift(resources []tfstate.Resource, live []aws.S3Bucket) []DriftResult {
	livemap := make(map[string]aws.S3Bucket, len(live))
	for _, b := range live {
		livemap[b.ID] = b
	}

	var results []DriftResult
	for _, r := range resources {
		if r.Type != "aws_s3_bucket" {
			continue
		}
		bucketID, ok := r.Attributes["bucket"].(string)
		if !ok || bucketID == "" {
			bucketID = r.Name
		}

		b, found := livemap[bucketID]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   bucketID,
				ResourceType: "aws_s3_bucket",
				Status:       StatusMissing,
				Details:      fmt.Sprintf("bucket %q not found in AWS", bucketID),
			})
			continue
		}

		expectedRegion, _ := r.Attributes["region"].(string)
		if expectedRegion != "" && b.Region != expectedRegion {
			results = append(results, DriftResult{
				ResourceID:   bucketID,
				ResourceType: "aws_s3_bucket",
				Status:       StatusDrifted,
				Details:      fmt.Sprintf("region mismatch: state=%s live=%s", expectedRegion, b.Region),
			})
		}
	}
	return results
}
