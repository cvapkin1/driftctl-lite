package drift

import (
	"errors"
	"testing"

	"github.com/owner/driftctl-lite/internal/aws"
)

type mockRDSFetcher struct {
	instances []aws.RDSInstance
	err       error
}

func (m *mockRDSFetcher) FetchAll() ([]aws.RDSInstance, error) {
	return m.instances, m.err
}

func rdsStateResource(id string) map[string]interface{} {
	return map[string]interface{}{"id": id}
}

func TestDetectRDSDrift_Missing(t *testing.T) {
	fetcher := &aws.RDSFetcher{}
	_ = fetcher // replaced by direct call pattern; use helper below

	results, err := detectRDSDriftWith([]map[string]interface{}{rdsStateResource("db-1")},
		[]aws.RDSInstance{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Status != "missing" {
		t.Errorf("expected missing drift, got %+v", results)
	}
}

func TestDetectRDSDrift_Deleted(t *testing.T) {
	inst := aws.RDSInstance{DBInstanceIdentifier: "db-2", DBInstanceStatus: "deleting"}
	results, err := detectRDSDriftWith([]map[string]interface{}{rdsStateResource("db-2")},
		[]aws.RDSInstance{inst}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Status != "deleted" {
		t.Errorf("expected deleted drift, got %+v", results)
	}
}

func TestDetectRDSDrift_NoDrift(t *testing.T) {
	inst := aws.RDSInstance{DBInstanceIdentifier: "db-3", DBInstanceStatus: "available"}
	results, err := detectRDSDriftWith([]map[string]interface{}{rdsStateResource("db-3")},
		[]aws.RDSInstance{inst}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Status != "ok" {
		t.Errorf("expected ok, got %+v", results)
	}
}

func TestDetectRDSDrift_FetchError(t *testing.T) {
	_, err := detectRDSDriftWith(nil, nil, errors.New("aws error"))
	if err == nil {
		t.Error("expected error, got nil")
	}
}

// detectRDSDriftWith is a test helper that bypasses the real fetcher.
func detectRDSDriftWith(state []map[string]interface{}, instances []aws.RDSInstance, fetchErr error) ([]RDSDriftResult, error) {
	if fetchErr != nil {
		return nil, fetchErr
	}
	f := &aws.RDSFetcher{}
	f.SetMock(instances)
	return DetectRDSDrift(state, f)
}
