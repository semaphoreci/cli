package cmd_get

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/semaphoreci/cli/cmd/handler"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"

	client "github.com/semaphoreci/cli/api/client"
)

var GetDashboardCmd = &cobra.Command{
	Use:     "dashboard [name]",
	Short:   "Get dashboards.",
	Long:    ``,
	Aliases: []string{"dashboards", "dash"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			c := client.DashboardsV1AlphaApi()

			dashList, err := c.ListDashboards()

			utils.Check(err)

			const padding = 3
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

			fmt.Fprintln(w, "NAME\tAGE")

			for _, d := range dashList {
				fmt.Fprintf(w, "%s\t%s\n", d.Metadata.Name, handler.RelativeAgeForHumans(d.UpdateTime))
			}

			w.Flush()
		} else {
			name := args[0]

			c := client.DashboardsV1AlphaApi()
			c.GetDashboard(name)

			params := semaphore_dashboards_v1alpha_dashboards_api.NewGetDashboardParams().WithIDOrName(name)
			resp, err := c.SemaphoreDashboardsV1alphaDashboardsAPI.GetDashboard(params)

			fmt.Printf("here")

			if resp, ok := err.(apiError); ok {
				fmt.Fprintf(os.Stderr, "error: (status %d) xx", resp.Code())

				return
			}

			utils.Check(err)

			j, err := resp.Payload.MarshalYaml()

			utils.Check(err)

			fmt.Printf("apiVersion: v1alpha\n")
			fmt.Printf("kind: Dashboard\n")
			fmt.Printf("%s", j)
		}
	},
}
