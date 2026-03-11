package resources

import (
	"github.com/prenansantana/mgc-infra-costs/internal/parser"
	"github.com/prenansantana/mgc-infra-costs/internal/pricing"
)

type ObjectStorageHandler struct{}

func (h *ObjectStorageHandler) ResourceType() string {
	return "mgc_object_storage_buckets"
}

func (h *ObjectStorageHandler) Estimate(resource parser.MGCResource, catalog *pricing.Catalog) (*ResourceEstimate, error) {
	// Object storage is usage-based (per GB stored + per GB transferred).
	// At creation time, the bucket has no data, so cost is R$ 0.00.
	// We show the per-GB price for reference.
	sku, _ := catalog.LookupByFlavor("Object_Storage", "cloud-std-os")

	pricePerGB := ""
	if sku != nil {
		pricePerGB = sku.MonthlyPrice
	}

	return &ResourceEstimate{
		Address:      resource.Address,
		ResourceType: resource.Type,
		MonthlyCost:  0,
		HourlyCost:   0,
		CostComponents: []CostComponent{
			{
				Name:            "Object Storage Standard (usage-based, R$ " + pricePerGB + "/GB/month)",
				Unit:            "GB/month",
				MonthlyQuantity: "0",
				Price:           pricePerGB,
				MonthlyCost:     0,
			},
		},
	}, nil
}
