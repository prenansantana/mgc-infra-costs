package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/prenansantana/mgc-infra-costs/internal/resources"
)

// InfracostOutput represents the Infracost-compatible JSON output format.
type InfracostOutput struct {
	Version  string             `json:"version"`
	Currency string             `json:"currency"`
	Projects []InfracostProject `json:"projects"`
}

type InfracostProject struct {
	Name      string             `json:"name"`
	Breakdown InfracostBreakdown `json:"breakdown"`
}

type InfracostBreakdown struct {
	Resources        []InfracostResource `json:"resources"`
	TotalMonthlyCost string              `json:"totalMonthlyCost"`
	TotalHourlyCost  string              `json:"totalHourlyCost"`
}

type InfracostResource struct {
	Name           string                   `json:"name"`
	ResourceType   string                   `json:"resourceType"`
	MonthlyCost    string                   `json:"monthlyCost"`
	HourlyCost     string                   `json:"hourlyCost,omitempty"`
	CostComponents []InfracostCostComponent `json:"costComponents"`
}

type InfracostCostComponent struct {
	Name            string `json:"name"`
	Unit            string `json:"unit"`
	HourlyQuantity  string `json:"hourlyQuantity,omitempty"`
	MonthlyQuantity string `json:"monthlyQuantity,omitempty"`
	Price           string `json:"price"`
	MonthlyCost     string `json:"monthlyCost"`
	HourlyCost      string `json:"hourlyCost,omitempty"`
}

func RenderInfracostJSON(w io.Writer, estimates []resources.ResourceEstimate, projectName string, totalMonthlyCost, totalHourlyCost float64) error {
	infracostResources := make([]InfracostResource, 0, len(estimates))

	for _, est := range estimates {
		components := make([]InfracostCostComponent, 0, len(est.CostComponents))
		for _, cc := range est.CostComponents {
			components = append(components, InfracostCostComponent{
				Name:            cc.Name,
				Unit:            cc.Unit,
				HourlyQuantity:  cc.HourlyQuantity,
				MonthlyQuantity: cc.MonthlyQuantity,
				Price:           cc.Price,
				MonthlyCost:     fmt.Sprintf("%.6f", cc.MonthlyCost),
				HourlyCost:      fmt.Sprintf("%.6f", cc.HourlyCost),
			})
		}

		infracostResources = append(infracostResources, InfracostResource{
			Name:           est.Address,
			ResourceType:   est.ResourceType,
			MonthlyCost:    fmt.Sprintf("%.6f", est.MonthlyCost),
			HourlyCost:     fmt.Sprintf("%.6f", est.HourlyCost),
			CostComponents: components,
		})
	}

	output := InfracostOutput{
		Version:  "0.2",
		Currency: "BRL",
		Projects: []InfracostProject{
			{
				Name: projectName,
				Breakdown: InfracostBreakdown{
					Resources:        infracostResources,
					TotalMonthlyCost: fmt.Sprintf("%.6f", totalMonthlyCost),
					TotalHourlyCost:  fmt.Sprintf("%.6f", totalHourlyCost),
				},
			},
		},
	}

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
