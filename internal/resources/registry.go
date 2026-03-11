package resources

import (
	"github.com/prenansantana/mgc-infra-costs/internal/parser"
	"github.com/prenansantana/mgc-infra-costs/internal/pricing"
)

// CostComponent represents a single line-item cost.
type CostComponent struct {
	Name            string  `json:"name"`
	Unit            string  `json:"unit"`
	HourlyQuantity  string  `json:"hourlyQuantity,omitempty"`
	MonthlyQuantity string  `json:"monthlyQuantity,omitempty"`
	Price           string  `json:"price"`
	MonthlyCost     float64 `json:"monthlyCost"`
	HourlyCost      float64 `json:"hourlyCost,omitempty"`
}

// ResourceEstimate represents the cost estimate for a single resource.
type ResourceEstimate struct {
	Address        string          `json:"name"`
	ResourceType   string          `json:"resourceType"`
	CostComponents []CostComponent `json:"costComponents"`
	MonthlyCost    float64         `json:"monthlyCost"`
	HourlyCost     float64         `json:"hourlyCost"`
}

// ResourceHandler calculates cost for a specific resource type.
type ResourceHandler interface {
	ResourceType() string
	Estimate(resource parser.MGCResource, catalog *pricing.Catalog) (*ResourceEstimate, error)
}

var handlers = map[string]ResourceHandler{}

func Register(h ResourceHandler) {
	handlers[h.ResourceType()] = h
}

func GetHandler(resourceType string) (ResourceHandler, bool) {
	h, ok := handlers[resourceType]
	return h, ok
}

func init() {
	Register(&VirtualMachineHandler{})
	Register(&BlockStorageHandler{})
	Register(&DatabaseHandler{})
	Register(&NetworkHandler{})
	Register(&ObjectStorageHandler{})
}
