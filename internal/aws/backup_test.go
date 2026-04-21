package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/backup"
	"github.com/stretchr/testify/assert"
)

func TestMapBackupVault_Fields(t *testing.T) {
	input := backup.BackupVaultListMember{
		BackupVaultName:  aws.String("my-vault"),
		BackupVaultArn:   aws.String("arn:aws:backup:us-east-1:123456789012:backup-vault:my-vault"),
		EncryptionKeyArn: aws.String("arn:aws:kms:us-east-1:123456789012:key/abc123"),
	}

	result := mapBackupVault(input)

	assert.Equal(t, "my-vault", result.Name)
	assert.Equal(t, "arn:aws:backup:us-east-1:123456789012:backup-vault:my-vault", result.ARN)
	assert.Equal(t, "arn:aws:kms:us-east-1:123456789012:key/abc123", result.EncryptionKeyARN)
}

func TestMapBackupVault_EmptyFields(t *testing.T) {
	input := backup.BackupVaultListMember{}

	result := mapBackupVault(input)

	assert.Equal(t, "", result.Name)
	assert.Equal(t, "", result.ARN)
	assert.Equal(t, "", result.EncryptionKeyARN)
}
