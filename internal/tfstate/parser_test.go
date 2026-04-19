package tfstate

import (
	"os"
	"testing"
)

const validState = `{
  "version": 4,
  "resources": [
    {
      "type": "aws_s3_bucket",
      "name": "my_bucket",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "attributes": {
            "id": "my-bucket-123",
            "region": "us-east-1"
          }
        }
      ]
    }
  ]
}`

func writeTempState(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "terraform-*.tfstate")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParseStateFile_Valid(t *testing.T) {
	path := writeTempState(t, validState)
	resources, err := ParseStateFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}
	r := resources[0]
	if r.Type != "aws_s3_bucket" {
		t.Errorf("expected type aws_s3_bucket, got %s", r.Type)
	}
	if r.Attributes["id"] != "my-bucket-123" {
		t.Errorf("unexpected id attribute: %v", r.Attributes["id"])
	}
}

func TestParseStateFile_UnsupportedVersion(t *testing.T) {
	path := writeTempState(t, `{"version": 3, "resources": []}`)
	_, err := ParseStateFile(path)
	if err == nil {
		t.Fatal("expected error for unsupported version, got nil")
	}
}

func TestParseStateFile_MissingFile(t *testing.T) {
	_, err := ParseStateFile("/nonexistent/path/terraform.tfstate")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
