package main

import (
	"fmt"
	"os"

	"github.com/prenansantana/mgc-infra-costs/internal/output"
	"github.com/prenansantana/mgc-infra-costs/internal/parser"
	"github.com/prenansantana/mgc-infra-costs/internal/pricing"
	"github.com/prenansantana/mgc-infra-costs/internal/resources"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "mgc-infra-costs",
		Short:   "Estimate Magalu Cloud infrastructure costs from Terraform plans",
		Version: version,
	}

	breakdownCmd := &cobra.Command{
		Use:   "breakdown",
		Short: "Show cost breakdown for a Terraform plan",
		RunE:  runBreakdown,
	}

	breakdownCmd.Flags().String("plan", "", "Path to terraform show -json output file")
	breakdownCmd.Flags().String("region", "br-se1", "MGC region for pricing")
	breakdownCmd.Flags().String("format", "table", "Output format: table, json")
	breakdownCmd.Flags().String("project-name", "", "Project name for JSON output (default: plan filename)")
	breakdownCmd.MarkFlagRequired("plan")

	rootCmd.AddCommand(breakdownCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runBreakdown(cmd *cobra.Command, args []string) error {
	planPath, _ := cmd.Flags().GetString("plan")
	region, _ := cmd.Flags().GetString("region")
	format, _ := cmd.Flags().GetString("format")
	projectName, _ := cmd.Flags().GetString("project-name")

	if projectName == "" {
		projectName = planPath
	}

	// Parse the Terraform plan
	plan, err := parser.ParsePlanFile(planPath)
	if err != nil {
		return fmt.Errorf("parsing plan: %w", err)
	}

	// Fetch pricing from MGC API
	fmt.Fprintf(os.Stderr, "Fetching pricing for region %s...\n", region)
	client := pricing.NewClient()
	categories, err := client.FetchPricing(region)
	if err != nil {
		return fmt.Errorf("fetching pricing: %w", err)
	}

	catalog := pricing.NewCatalog(categories)

	// Extract MGC resources
	mgcResources := parser.ExtractMGCResources(plan)
	if len(mgcResources) == 0 {
		fmt.Fprintln(os.Stderr, "No MGC resources found in the plan.")
		return nil
	}

	fmt.Fprintf(os.Stderr, "Found %d MGC resource(s).\n", len(mgcResources))

	// Estimate costs
	var estimates []resources.ResourceEstimate
	var totalMonthlyCost, totalHourlyCost float64
	var warnings []string

	for _, r := range mgcResources {
		handler, ok := resources.GetHandler(r.Type)
		if !ok {
			warnings = append(warnings, fmt.Sprintf("No cost handler for resource type %s (%s)", r.Type, r.Address))
			continue
		}

		estimate, err := handler.Estimate(r, catalog)
		if err != nil {
			warnings = append(warnings, fmt.Sprintf("Could not estimate %s: %v", r.Address, err))
			continue
		}

		estimates = append(estimates, *estimate)
		totalMonthlyCost += estimate.MonthlyCost
		totalHourlyCost += estimate.HourlyCost
	}

	// Print warnings
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", w)
	}

	// Render output
	switch format {
	case "json":
		return output.RenderInfracostJSON(os.Stdout, estimates, projectName, totalMonthlyCost, totalHourlyCost)
	case "table":
		output.RenderTable(os.Stdout, estimates, totalMonthlyCost)
	default:
		return fmt.Errorf("unknown format: %s", format)
	}

	return nil
}
