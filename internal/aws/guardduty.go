package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
)

type GuardDutyDetector struct {
	ID     string
	Status string
	ARN    string
}

type GuardDutyFetcher struct {
	client *guardduty.Client
}

func NewGuardDutyFetcher(client *guardduty.Client) *GuardDutyFetcher {
	return &GuardDutyFetcher{client: client}
}

func (f *GuardDutyFetcher) FetchAll(ctx context.Context) ([]GuardDutyDetector, error) {
	listOut, err := f.client.ListDetectors(ctx, &guardduty.ListDetectorsInput{})
	if err != nil {
		return nil, err
	}

	var detectors []GuardDutyDetector
	for _, id := range listOut.DetectorIds {
		getOut, err := f.client.GetDetector(ctx, &guardduty.GetDetectorInput{
			DetectorId: aws.String(id),
		})
		if err != nil {
			return nil, err
		}
		detectors = append(detectors, mapGuardDutyDetector(id, getOut))
	}
	return detectors, nil
}

func mapGuardDutyDetector(id string, out *guardduty.GetDetectorOutput) GuardDutyDetector {
	status := ""
	if out.Status != "" {
		status = string(out.Status)
	}
	arn := ""
	if out.ServiceRole != nil {
		arn = *out.ServiceRole
	}
	return GuardDutyDetector{
		ID:     id,
		Status: status,
		ARN:    arn,
	}
}
