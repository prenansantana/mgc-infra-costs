package resources

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/prenansantana/mgc-infra-costs/internal/parser"
	"github.com/prenansantana/mgc-infra-costs/internal/pricing"
)

// shortFormatRegex matches machine types like "BV2-4-20", "DP4-8-100"
var shortFormatRegex = regexp.MustCompile(`^([A-Z]+)(\d+)-(\d+)-(\d+)$`)

type VirtualMachineHandler struct{}

func (h *VirtualMachineHandler) ResourceType() string {
	return "mgc_virtual_machine_instances"
}

func (h *VirtualMachineHandler) Estimate(resource parser.MGCResource, catalog *pricing.Catalog) (*ResourceEstimate, error) {
	machineType := extractMachineType(resource.Attributes)
	if machineType == "" {
		return nil, fmt.Errorf("machine_type not found for %s", resource.Address)
	}

	sku, err := lookupVMSKU(catalog, machineType)
	if err != nil {
		return nil, fmt.Errorf("pricing for %s (machine_type=%s): %w", resource.Address, machineType, err)
	}

	monthlyPrice := pricing.MonthlyPriceFloat(sku)
	hourlyPrice := pricing.HourlyPriceFloat(sku)

	desc := fmt.Sprintf("Instance (%s, %s)", machineType, sku.SKUClass)
	if sku.VCPUs > 0 && sku.RAM != nil && sku.Disk != nil {
		desc = fmt.Sprintf("Instance (%s, %s - %d vCPU, %d%s RAM, %d%s Disk)",
			machineType, sku.SKUClass, sku.VCPUs,
			sku.RAM.Value, sku.RAM.Unit,
			sku.Disk.Value, sku.Disk.Unit)
	}

	return &ResourceEstimate{
		Address:      resource.Address,
		ResourceType: resource.Type,
		MonthlyCost:  monthlyPrice,
		HourlyCost:   hourlyPrice,
		CostComponents: []CostComponent{
			{
				Name:            desc,
				Unit:            "hours",
				HourlyQuantity:  "1",
				MonthlyQuantity: "730",
				Price:           sku.HourlyPrice,
				MonthlyCost:     monthlyPrice,
				HourlyCost:      hourlyPrice,
			},
		},
	}, nil
}

// lookupVMSKU tries to find a VM SKU by:
// 1. Direct flavor_name match (e.g., "i1-c2-r4-d20")
// 2. Short format parsing (e.g., "BV2-4-20" → class=BV, vcpu=2, ram=4, disk=20)
func lookupVMSKU(catalog *pricing.Catalog, machineType string) (*pricing.SKU, error) {
	// Try direct match first
	sku, err := catalog.LookupByFlavor("Virtual_Machines", machineType)
	if err == nil {
		return sku, nil
	}

	// Try parsing short format: BV2-4-20, DP4-8-100, etc.
	matches := shortFormatRegex.FindStringSubmatch(machineType)
	if matches == nil {
		return nil, fmt.Errorf("unknown machine_type format %q", machineType)
	}

	classPrefix := matches[1]
	vcpu, _ := strconv.Atoi(matches[2])
	ram, _ := strconv.Atoi(matches[3])
	disk, _ := strconv.Atoi(matches[4])

	return catalog.LookupBySpecs("Virtual_Machines", classPrefix, vcpu, ram, disk)
}

func extractMachineType(attrs map[string]interface{}) string {
	if mt, ok := attrs["machine_type"]; ok {
		switch v := mt.(type) {
		case map[string]interface{}:
			if name, ok := v["name"].(string); ok {
				return name
			}
		case string:
			return v
		}
	}
	if mt, ok := attrs["machine_type_id"].(string); ok {
		return mt
	}
	return ""
}
