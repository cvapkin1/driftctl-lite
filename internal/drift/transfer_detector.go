package drift

import (
	"fmt"

	"github.com/jonmorehouse/driftctl-lite/internal/aws"
)

type TransferDriftResult struct {
	ServerID string
	Status   string
	Message  string
}

// DetectTransferDrift compares Terraform state resources against live AWS Transfer servers.
func DetectTransferDrift(stateResources []map[string]interface{}, live []aws.TransferServer) []TransferDriftResult {
	liveMap := make(map[string]aws.TransferServer)
	for _, s := range live {
		liveMap[s.ServerID] = s
	}

	var results []TransferDriftResult

	for _, res := range stateResources {
		serverID, _ := res["server_id"].(string)
		expectedARN, _ := res["arn"].(string)

		liveServer, found := liveMap[serverID]
		if !found {
			results = append(results, TransferDriftResult{
				ServerID: serverID,
				Status:   "missing",
				Message:  fmt.Sprintf("Transfer server %q not found in AWS", serverID),
			})
			continue
		}

		if liveServer.State == "OFFLINE" || liveServer.State == "STOP_FAILED" {
			results = append(results, TransferDriftResult{
				ServerID: serverID,
				Status:   "inactive",
				Message:  fmt.Sprintf("Transfer server %q is in state %q", serverID, liveServer.State),
			})
			continue
		}

		if expectedARN != "" && liveServer.ARN != expectedARN {
			results = append(results, TransferDriftResult{
				ServerID: serverID,
				Status:   "arn_mismatch",
				Message:  fmt.Sprintf("Transfer server %q ARN mismatch: state=%q live=%q", serverID, expectedARN, liveServer.ARN),
			})
			continue
		}

		results = append(results, TransferDriftResult{
			ServerID: serverID,
			Status:   "ok",
			Message:  fmt.Sprintf("Transfer server %q is in sync", serverID),
		})
	}

	return results
}
