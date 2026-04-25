package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apprunner"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AppRunnerService struct {
	ID     string
	Name   string
	ARN    string
	Status string
	URL    string
}

type AppRunnerFetcher struct {
	client *apprunner.Client
}

func NewAppRunnerFetcher(ctx context.Context, region string) (*AppRunnerFetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &AppRunnerFetcher{client: apprunner.NewFromConfig(cfg)}, nil
}

func (f *AppRunnerFetcher) FetchAll(ctx context.Context) ([]AppRunnerService, error) {
	var services []AppRunnerService
	var nextToken *string

	for {
		resp, err := f.client.ListServices(ctx, &apprunner.ListServicesInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, svc := range resp.ServiceSummaryList {
			services = append(services, mapAppRunnerService(svc))
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return services, nil
}

func mapAppRunnerService(svc apprunnerTypes.ServiceSummary) AppRunnerService {
	return AppRunnerService{
		ID:     aws.ToString(svc.ServiceId),
		Name:   aws.ToString(svc.ServiceName),
		ARN:    aws.ToString(svc.ServiceArn),
		Status: string(svc.Status),
		URL:    aws.ToString(svc.ServiceUrl),
	}
}
