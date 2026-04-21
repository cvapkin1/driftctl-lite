package drift

import (
	"fmt"

	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

const (
	acmResourceType    = "aws_acm_certificate"
	acmStatusInvalid   = "INACTIVE"
	acmStatusExpired   = "EXPIRED"
	acmStatusRevoked   = "REVOKED"
)

func DetectACMDrift(resources []tfstate.Resource, certs []aws.ACMCertificate) []DriftResult {
	var results []DriftResult

	live := make(map[string]aws.ACMCertificate)
	for _, c := range certs {
		live[c.ARN] = c
	}

	for _, res := range resources {
		if res.Type != acmResourceType {
			continue
		}
		arn, _ := res.Attributes["arn"].(string)
		if arn == "" {
			continue
		}

		cert, found := live[arn]
		if !found {
			results = append(results, DriftResult{
				ResourceID:   arn,
				ResourceType: acmResourceType,
				DriftType:    "missing",
				Details:      fmt.Sprintf("certificate %s not found in AWS", arn),
			})
			continue
		}

		if cert.Status == acmStatusExpired || cert.Status == acmStatusRevoked || cert.Status == acmStatusInvalid {
			results = append(results, DriftResult{
				ResourceID:   arn,
				ResourceType: acmResourceType,
				DriftType:    "status_drift",
				Details:      fmt.Sprintf("certificate %s has unexpected status: %s", arn, cert.Status),
			})
		}

		expectedDomain, _ := res.Attributes["domain_name"].(string)
		if expectedDomain != "" && cert.Domain != expectedDomain {
			results = append(results, DriftResult{
				ResourceID:   arn,
				ResourceType: acmResourceType,
				DriftType:    "domain_mismatch",
				Details:      fmt.Sprintf("certificate %s domain mismatch: state=%s live=%s", arn, expectedDomain, cert.Domain),
			})
		}
	}
	return results
}
