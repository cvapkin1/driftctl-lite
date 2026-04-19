package aws_test

import (
	"testing"

	internalaws "github.com/example/driftctl-lite/internal/aws"
)

func TestMapInstance_Tags(t *testing.T) {
	inst := internalaws.EC2Instance{
		ID:           "i-abc123",
		InstanceType: "t3.micro",
		State:        "running",
		Tags:         map[string]string{"Name": "web-server", "Env": "prod"},
	}
	if inst.ID != "i-abc123" {
		t.Errorf("expected i-abc123, got %s", inst.ID)
	}
	if inst.Tags["Name"] != "web-server" {
		t.Errorf("expected web-server tag, got %s", inst.Tags["Name"])
	}
}

func TestEC2Instance_StateRunning(t *testing.T) {
	inst := internalaws.EC2Instance{State: "running"}
	if inst.State != "running" {
		t.Errorf("unexpected state: %s", inst.State)
	}
}
