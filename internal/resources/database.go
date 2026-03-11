package resources

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/prenansantana/mgc-infra-costs/internal/parser"
	"github.com/prenansantana/mgc-infra-costs/internal/pricing"
)

// dbShortFormatRegex matches DBaaS instance types like "BV1-4-10"
var dbShortFormatRegex = regexp.MustCompile(`^([A-Z]+)(\d+)-(\d+)-(\d+)$`)

type DatabaseHandler struct{}

func (h *DatabaseHandler) ResourceType() string {
	return "mgc_dbaas_instances"
}

func (h *DatabaseHandler) Estimate(resource parser.MGCResource, catalog *pricing.Catalog) (*ResourceEstimate, error) {
	instanceType := extractDBInstanceType(resource.Attributes)
	if instanceType == "" {
		return nil, fmt.Errorf("instance_type/flavor not found for %s", resource.Address)
	}

	engineName := extractDBEngine(resource.Attributes)

	sku, err := lookupDBSKU(catalog, instanceType, engineName)
	if err != nil {
		return nil, fmt.Errorf("pricing for %s (instance_type=%s, engine=%s): %w", resource.Address, instanceType, engineName, err)
	}

	monthlyPrice := pricing.MonthlyPriceFloat(sku)
	hourlyPrice := pricing.HourlyPriceFloat(sku)

	desc := fmt.Sprintf("DBaaS %s (%s)", sku.System, sku.Description)
	if sku.Availability != "" {
		desc = fmt.Sprintf("DBaaS %s %s (%s)", sku.System, sku.Availability, sku.Description)
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

// lookupDBSKU tries to find a DBaaS SKU by:
// 1. Direct flavor_name match
// 2. Short format parsing (e.g., "BV1-4-10" → class=BV, vcpu=1, ram=4, disk=10)
// 3. Pattern matching
func lookupDBSKU(catalog *pricing.Catalog, instanceType, engineName string) (*pricing.SKU, error) {
	// Try direct match
	sku, err := catalog.LookupByFlavor("Dbaas", instanceType)
	if err == nil {
		return sku, nil
	}

	// Try short format: BV1-4-10, DP2-16-40
	matches := dbShortFormatRegex.FindStringSubmatch(instanceType)
	if matches != nil {
		classPrefix := matches[1]
		vcpu, _ := strconv.Atoi(matches[2])
		ram, _ := strconv.Atoi(matches[3])
		disk, _ := strconv.Atoi(matches[4])

		found := findDBSKUBySpecs(catalog, classPrefix, vcpu, ram, disk, engineName)
		if found != nil {
			return found, nil
		}
	}

	// Try pattern matching as fallback
	found := findDBSKUByPattern(catalog, instanceType)
	if found != nil {
		return found, nil
	}

	return nil, fmt.Errorf("no DBaaS SKU found for instance_type=%q engine=%q", instanceType, engineName)
}

func findDBSKUBySpecs(catalog *pricing.Catalog, classPrefix string, vcpu, ramGB, diskGB int, engineName string) *pricing.SKU {
	skus := catalog.ListFlavors("Dbaas")
	classTag := fmt.Sprintf("(%s)", classPrefix)

	for _, sku := range skus {
		if sku.VCPUs != vcpu || sku.RAM == nil || sku.Disk == nil {
			continue
		}
		if sku.RAM.Value != ramGB || sku.Disk.Value != diskGB {
			continue
		}
		if !strings.Contains(sku.SKUClass, classTag) {
			continue
		}
		// Match engine if specified
		if engineName != "" && !strings.EqualFold(sku.System, engineName) {
			continue
		}
		return &sku
	}

	// Retry without engine filter if no match
	if engineName != "" {
		for _, sku := range skus {
			if sku.VCPUs != vcpu || sku.RAM == nil || sku.Disk == nil {
				continue
			}
			if sku.RAM.Value != ramGB || sku.Disk.Value != diskGB {
				continue
			}
			if !strings.Contains(sku.SKUClass, classTag) {
				continue
			}
			return &sku
		}
	}
	return nil
}

func extractDBInstanceType(attrs map[string]interface{}) string {
	if it, ok := attrs["instance_type"]; ok {
		switch v := it.(type) {
		case map[string]interface{}:
			if name, ok := v["name"].(string); ok {
				return name
			}
		case string:
			return v
		}
	}
	if id, ok := attrs["instance_type_id"].(string); ok {
		return id
	}
	if fn, ok := attrs["flavor_name"].(string); ok {
		return fn
	}
	return ""
}

func extractDBEngine(attrs map[string]interface{}) string {
	if en, ok := attrs["engine_name"].(string); ok {
		return en
	}
	if en, ok := attrs["engine"].(string); ok {
		return en
	}
	return ""
}

func findDBSKUByPattern(catalog *pricing.Catalog, pattern string) *pricing.SKU {
	skus := catalog.ListFlavors("Dbaas")
	for _, sku := range skus {
		if strings.Contains(sku.FlavorName, pattern) || strings.Contains(pattern, sku.FlavorName) {
			return &sku
		}
	}
	return nil
}
