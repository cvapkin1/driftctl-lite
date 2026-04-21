package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
)

type BackupVault struct {
	Name             string
	ARN              string
	EncryptionKeyARN string
}

type BackupVaultFetcher struct {
	client *backup.Client
}

func NewBackupFetcher(client *backup.Client) *BackupVaultFetcher {
	return &BackupVaultFetcher{client: client}
}

func (f *BackupVaultFetcher) FetchAll(ctx context.Context) ([]BackupVault, error) {
	var vaults []BackupVault
	var nextToken *string

	for {
		resp, err := f.client.ListBackupVaults(ctx, &backup.ListBackupVaultsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, err
		}

		for _, v := range resp.BackupVaultList {
			vaults = append(vaults, mapBackupVault(v))
		}

		if resp.NextToken == nil {
			break
		}
		nextToken = resp.NextToken
	}

	return vaults, nil
}

func mapBackupVault(v backup.BackupVaultListMember) BackupVault {
	return BackupVault{
		Name:             aws.ToString(v.BackupVaultName),
		ARN:              aws.ToString(v.BackupVaultArn),
		EncryptionKeyARN: aws.ToString(v.EncryptionKeyArn),
	}
}
