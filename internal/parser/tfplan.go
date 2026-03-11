package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// TFPlan represents the terraform show -json output structure.
type TFPlan struct {
	PlannedValues   PlannedValues    `json:"planned_values"`
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

type PlannedValues struct {
	RootModule Module `json:"root_module"`
}

type Module struct {
	Resources    []Resource `json:"resources"`
	ChildModules []Module   `json:"child_modules"`
}

type Resource struct {
	Address      string                 `json:"address"`
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	ProviderName string                 `json:"provider_name"`
	Values       map[string]interface{} `json:"values"`
}

type ResourceChange struct {
	Address string `json:"address"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Change  Change `json:"change"`
}

type Change struct {
	Actions []string               `json:"actions"`
	Before  map[string]interface{} `json:"before"`
	After   map[string]interface{} `json:"after"`
}

// MGCResource represents a parsed MGC resource with its relevant attributes.
type MGCResource struct {
	Address    string
	Type       string
	Name       string
	Attributes map[string]interface{}
}

var supportedTypes = map[string]bool{
	"mgc_virtual_machine_instances": true,
	"mgc_block_storage_volumes":     true,
	"mgc_dbaas_instances":           true,
	"mgc_network_public_ips":        true,
	"mgc_object_storage_buckets":    true,
}

func ParsePlanFile(path string) (*TFPlan, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading plan file: %w", err)
	}
	return ParsePlan(data)
}

func ParsePlan(data []byte) (*TFPlan, error) {
	var plan TFPlan
	if err := json.Unmarshal(data, &plan); err != nil {
		return nil, fmt.Errorf("parsing plan JSON: %w", err)
	}
	return &plan, nil
}

// ExtractMGCResources extracts all MGC resources from the plan.
func ExtractMGCResources(plan *TFPlan) []MGCResource {
	var resources []MGCResource
	extractFromModule(&plan.PlannedValues.RootModule, &resources)
	return resources
}

func extractFromModule(mod *Module, resources *[]MGCResource) {
	for _, r := range mod.Resources {
		if !strings.HasPrefix(r.Type, "mgc_") {
			continue
		}
		if !supportedTypes[r.Type] {
			continue
		}
		*resources = append(*resources, MGCResource{
			Address:    r.Address,
			Type:       r.Type,
			Name:       r.Name,
			Attributes: r.Values,
		})
	}
	for i := range mod.ChildModules {
		extractFromModule(&mod.ChildModules[i], resources)
	}
}

// ExtractMGCResourceChanges extracts MGC resources from resource_changes (for diff).
func ExtractMGCResourceChanges(plan *TFPlan) []ResourceChange {
	var changes []ResourceChange
	for _, rc := range plan.ResourceChanges {
		if !strings.HasPrefix(rc.Type, "mgc_") {
			continue
		}
		if !supportedTypes[rc.Type] {
			continue
		}
		changes = append(changes, rc)
	}
	return changes
}
