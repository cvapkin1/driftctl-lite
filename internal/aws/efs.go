package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	efstypes "github.com/aws/aws-sdk-go-v2/service/efs/types"
)

type EFSFileSystem struct {
	ID             string
	Name           string
	LifeCycleState string
	Encrypted      bool
	ThroughputMode string
}

type EFSClient interface {
	DescribeFileSystems(ctx context.Context, params *efs.DescribeFileSystemsInput, optFns ...func(*efs.Options)) (*efs.DescribeFileSystemsOutput, error)
}

type EFSFetcher struct {
	client EFSClient
}

func NewEFSFetcher(ctx context.Context, region string) (*EFSFetcher, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return &EFSFetcher{client: efs.NewFromConfig(cfg)}, nil
}

func (f *EFSFetcher) FetchAll(ctx context.Context) ([]EFSFileSystem, error) {
	var results []EFSFileSystem
	var marker *string
	for {
		out, err := f.client.DescribeFileSystems(ctx, &efs.DescribeFileSystemsInput{
			Marker: marker,
		})
		if err != nil {
			return nil, err
		}
		for _, fs := range out.FileSystems {
			results = append(results, mapEFSFileSystem(fs))
		}
		if out.NextMarker == nil {
			break
		}
		marker = out.NextMarker
	}
	return results, nil
}

func mapEFSFileSystem(fs efstypes.FileSystemDescription) EFSFileSystem {
	name := ""
	if fs.Name != nil {
		name = aws.ToString(fs.Name)
	}
	return EFSFileSystem{
		ID:             aws.ToString(fs.FileSystemId),
		Name:           name,
		LifeCycleState: string(fs.LifeCycleState),
		Encrypted:      aws.ToBool(fs.Encrypted),
		ThroughputMode: string(fs.ThroughputMode),
	}
}
