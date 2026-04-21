package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func TestMapECSCluster_Fields(t *testing.T) {
	c := ecs.Cluster{
		ClusterArn:  aws.String("arn:aws:ecs:us-east-1:123456789012:cluster/my-cluster"),
		ClusterName: aws.String("my-cluster"),
		Status:      aws.String("ACTIVE"),
	}

	result := mapECSCluster(c)

	if result.ARN != "arn:aws:ecs:us-east-1:123456789012:cluster/my-cluster" {
		t.Errorf("expected ARN to match, got %s", result.ARN)
	}
	if result.Name != "my-cluster" {
		t.Errorf("expected Name to be my-cluster, got %s", result.Name)
	}
	if result.Status != "ACTIVE" {
		t.Errorf("expected Status to be ACTIVE, got %s", result.Status)
	}
}

func TestMapECSCluster_EmptyFields(t *testing.T) {
	c := ecs.Cluster{}
	result := mapECSCluster(c)

	if result.ARN != "" {
		t.Errorf("expected empty ARN, got %s", result.ARN)
	}
	if result.Status != "" {
		t.Errorf("expected empty Status, got %s", result.Status)
	}
}
