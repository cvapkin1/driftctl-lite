package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/transfer"
	"github.com/aws/aws-sdk-go-v2/service/transfer/types"
)

type TransferServer struct {
	ServerID string
	ARN      string
	State    string
	Endpoint string
	Protocols []string
}

type TransferFetcher struct {
	client *transfer.Client
}

func NewTransferFetcher(cfg aws.Config) *TransferFetcher {
	return &TransferFetcher{
		client: transfer.NewFromConfig(cfg),
	}
}

func (f *TransferFetcher) FetchAll(ctx context.Context) ([]TransferServer, error) {
	var servers []TransferServer
	paginator := transfer.NewListServersPaginator(f.client, &transfer.ListServersInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, s := range page.Servers {
			servers = append(servers, mapTransferServer(s))
		}
	}
	return servers, nil
}

func mapTransferServer(s types.ListedServer) TransferServer {
	var protocols []string
	for _, p := range s.Protocols {
		protocols = append(protocols, string(p))
	}
	endpoint := ""
	if s.EndpointType != "" {
		endpoint = string(s.EndpointType)
	}
	return TransferServer{
		ServerID:  aws.ToString(s.ServerId),
		ARN:       aws.ToString(s.Arn),
		State:     string(s.State),
		Endpoint:  endpoint,
		Protocols: protocols,
	}
}
