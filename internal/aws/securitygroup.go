package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type SecurityGroup struct {
	ID          string
	Name        string
	Description string
	VPCID       string
	Tags        map[string]string
}

type SecurityGroupFetcher struct {
	client *ec2.Client
}

func NewSecurityGroupFetcher(client *ec2.Client) *SecurityGroupFetcher {
	return &SecurityGroupFetcher{client: client}
}

func (f *SecurityGroupFetcher) FetchAll(ctx context.Context) ([]SecurityGroup, error) {
	var groups []SecurityGroup
	paginator := ec2.NewDescribeSecurityGroupsPaginator(f.client, &ec2.DescribeSecurityGroupsInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}
		for _, sg := range page.SecurityGroups {
			groups = append(groups, mapSecurityGroup(sg))
		}
	}
	return groups, nil
}

func mapSecurityGroup(sg types.SecurityGroup) SecurityGroup {
	tags := make(map[string]string)
	for _, t := range sg.Tags {
		if t.Key != nil && t.Value != nil {
			tags[aws.ToString(t.Key)] = aws.ToString(t.Value)
		}
	}
	return SecurityGroup{
		ID:          aws.ToString(sg.GroupId),
		Name:        aws.ToString(sg.GroupName),
		Description: aws.ToString(sg.Description),
		VPCID:       aws.ToString(sg.VpcId),
		Tags:        tags,
	}
}
