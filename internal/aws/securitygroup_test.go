package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
)

func TestMapSecurityGroup_Fields(t *testing.T) {
	sg := types.SecurityGroup{
		GroupId:     aws.String("sg-0abc123"),
		GroupName:   aws.String("my-sg"),
		Description: aws.String("test security group"),
		VpcId:       aws.String("vpc-0123456"),
		Tags: []types.Tag{
			{Key: aws.String("Env"), Value: aws.String("prod")},
		},
	}

	result := mapSecurityGroup(sg)

	assert.Equal(t, "sg-0abc123", result.ID)
	assert.Equal(t, "my-sg", result.Name)
	assert.Equal(t, "test security group", result.Description)
	assert.Equal(t, "vpc-0123456", result.VPCID)
	assert.Equal(t, "prod", result.Tags["Env"])
}

func TestMapSecurityGroup_EmptyFields(t *testing.T) {
	sg := types.SecurityGroup{
		GroupId:   aws.String("sg-empty"),
		GroupName: aws.String("default"),
	}

	result := mapSecurityGroup(sg)

	assert.Equal(t, "sg-empty", result.ID)
	assert.Equal(t, "default", result.Name)
	assert.Empty(t, result.Description)
	assert.Empty(t, result.VPCID)
	assert.Empty(t, result.Tags)
}
