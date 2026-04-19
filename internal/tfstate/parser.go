package tfstate

import (
	"encoding/json"
	"fmt"
	"os"
)

// Resource represents a single resource extracted from Terraform state.
type Resource struct {
	Type       string            `json:"type"`
	Name       string            `json:"name"`
	Provider   string            `json:"provider"`
	Attributes map[string]interface{} `json:"attributes"`
}

type tfStateFile struct {
	Version   int           `json:"version"`
	Resources []tfResource  `json:"resources"`
}

type tfResource struct {
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	Provider  string       `json:"provider"`
	Instances []tfInstance `json:"instances"`
}

type tfInstance struct {
	Attributes map[string]interface{} `json:"attributes"`
}

// ParseStateFile reads and parses a Terraform state file at the given path.
// It returns a flat list of Resource structs.
func ParseStateFile(path string) ([]Resource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading state file: %w", err)
	}

	var state tfStateFile
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("parsing state file: %w", err)
	}

	if state.Version < 4 {
		return nil, fmt.Errorf("unsupported state version %d (minimum: 4)", state.Version)
	}

	var resources []Resource
	for _, r := range state.Resources {
		for _, inst := range r.Instances {
			resources = append(resources, Resource{
				Type:       r.Type,
				Name:       r.Name,
				Provider:   r.Provider,
				Attributes: inst.Attributes,
			})
		}
	}
	return resources, nil
}
