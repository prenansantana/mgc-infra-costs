package resources

import (
	"github.com/prenansantana/mgc-infra-costs/internal/parser"
	"github.com/prenansantana/mgc-infra-costs/internal/pricing"
)

type NetworkHandler struct{}

func (h *NetworkHandler) ResourceType() string {
	return "mgc_network_public_ips"
}

func (h *NetworkHandler) Estimate(resource parser.MGCResource, catalog *pricing.Catalog) (*ResourceEstimate, error) {
	// Public IPs in MGC have no allocation cost.
	// Cost is only incurred on data transfer egress, which is usage-based.
	_ = catalog

	return &ResourceEstimate{
		Address:      resource.Address,
		ResourceType: resource.Type,
		MonthlyCost:  0,
		HourlyCost:   0,
		CostComponents: []CostComponent{
			{
				Name:        "Public IP allocation (no charge)",
				Unit:        "hours",
				MonthlyCost: 0,
			},
		},
	}, nil
}
