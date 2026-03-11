package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/prenansantana/mgc-infra-costs/internal/resources"
)

func RenderTable(w io.Writer, estimates []resources.ResourceEstimate, totalMonthlyCost float64) {
	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)

	fmt.Fprintln(tw, "")
	fmt.Fprintln(tw, " Name\tMonthly Cost\t")
	fmt.Fprintln(tw, " "+strings.Repeat("─", 60)+"\t"+strings.Repeat("─", 15)+"\t")

	for _, est := range estimates {
		fmt.Fprintf(tw, " %s\t\t\n", est.Address)
		for _, cc := range est.CostComponents {
			fmt.Fprintf(tw, "   └─ %s\tR$ %.2f\t\n", cc.Name, cc.MonthlyCost)
		}
	}

	fmt.Fprintln(tw, " "+strings.Repeat("─", 60)+"\t"+strings.Repeat("─", 15)+"\t")
	fmt.Fprintf(tw, " TOTAL MONTHLY COST\tR$ %.2f\t\n", totalMonthlyCost)
	fmt.Fprintln(tw, "")

	tw.Flush()
}
