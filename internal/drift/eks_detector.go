package drift

import (
	"fmt"

	"github.com/your-org/driftctl-lite/internal/aws"
	"github.com/your-org/driftctl-lite/internal/tfstate"
)

// DetectEKSDrift compares EKS clusters in Terraform state against live AWS clusters.
func DetectEKSDrift(resources []tfstate.Resource, live []aws.EKSCluster) []string {
	var drifts []string

	liveMap := make(map[string]aws.EKSCluster)
	for _, c := range live {
		liveMap[c.Name] = c
	}

	for _, r := range resources {
		if r.Type != "aws_eks_cluster" {
			continue
		}

		name, _ := r.Attributes["name"].(string)
		if name == "" {
			continue
		}

		liveCluster, found := liveMap[name]
		if !found {
			drifts = append(drifts, fmt.Sprintf("EKS cluster %q is in state but not found in AWS", name))
			continue
		}

		if liveCluster.Status == "DELETING" || liveCluster.Status == "FAILED" {
			drifts = append(drifts, fmt.Sprintf("EKS cluster %q has unexpected status: %s", name, liveCluster.Status))
		}

		if stateVersion, ok := r.Attributes["version"].(string); ok && stateVersion != "" {
			if liveCluster.Version != stateVersion {
				drifts = append(drifts, fmt.Sprintf("EKS cluster %q version mismatch: state=%s live=%s", name, stateVersion, liveCluster.Version))
			}
		}
	}

	return drifts
}
