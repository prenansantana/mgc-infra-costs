package resources

import (
	"fmt"

	"github.com/prenansantana/mgc-infra-costs/internal/parser"
	"github.com/prenansantana/mgc-infra-costs/internal/pricing"
)

type BlockStorageHandler struct{}

func (h *BlockStorageHandler) ResourceType() string {
	return "mgc_block_storage_volumes"
}

func (h *BlockStorageHandler) Estimate(resource parser.MGCResource, catalog *pricing.Catalog) (*ResourceEstimate, error) {
	volumeType := extractVolumeType(resource.Attributes)
	if volumeType == "" {
		volumeType = "cloud_nvme" // default
	}

	sizeGB := extractVolumeSize(resource.Attributes)
	if sizeGB == 0 {
		sizeGB = 1 // minimum
	}

	sku, err := catalog.LookupByFlavor("Block_Storage", volumeType)
	if err != nil {
		return nil, fmt.Errorf("pricing for %s (type=%s): %w", resource.Address, volumeType, err)
	}

	pricePerGB := pricing.MonthlyPriceFloat(sku)
	monthlyCost := pricePerGB * float64(sizeGB)

	return &ResourceEstimate{
		Address:      resource.Address,
		ResourceType: resource.Type,
		MonthlyCost:  monthlyCost,
		CostComponents: []CostComponent{
			{
				Name:            fmt.Sprintf("Storage (%s, %dGB)", volumeType, sizeGB),
				Unit:            "GB/month",
				MonthlyQuantity: fmt.Sprintf("%d", sizeGB),
				Price:           sku.MonthlyPrice,
				MonthlyCost:     monthlyCost,
			},
		},
	}, nil
}

func extractVolumeType(attrs map[string]interface{}) string {
	if t, ok := attrs["type"]; ok {
		switch v := t.(type) {
		case map[string]interface{}:
			if name, ok := v["name"].(string); ok {
				return name
			}
		case string:
			return v
		}
	}
	if t, ok := attrs["type_name"].(string); ok {
		return t
	}
	return ""
}

func extractVolumeSize(attrs map[string]interface{}) int {
	if size, ok := attrs["size"].(float64); ok {
		return int(size)
	}
	if disk, ok := attrs["disk"].(map[string]interface{}); ok {
		if size, ok := disk["size"].(float64); ok {
			return int(size)
		}
	}
	return 0
}
