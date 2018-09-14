package cmd_get

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/semaphoreci/cli/api"
	"github.com/semaphoreci/cli/api/client/semaphore_dashboards_v1alpha_dashboards_api"
	"github.com/semaphoreci/cli/cmd/handler"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var GetDashboardCmd = &cobra.Command{
	Use:     "dashboard [name]",
	Short:   "Get dashboards.",
	Long:    ``,
	Aliases: []string{"dashboards", "dash"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			c := api.DefaultClient()

			resp, err := c.SemaphoreDashboardsV1alphaDashboardsAPI.ListDashboards(nil)

			utils.Check(err)

			const padding = 3
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

			fmt.Fprintln(w, "NAME\tAGE")

			for _, d := range resp.Payload.Dashboards {
				update_time, err := strconv.ParseInt(d.Metadata.UpdateTime, 10, 64)

				utils.Check(err)

				fmt.Fprintf(w, "%s\t%s\n", d.Metadata.Name, handler.RelativeAgeForHumans(update_time))
			}

			w.Flush()
		} else {
			name := args[0]

			c := api.DefaultClient()

			params := semaphore_dashboards_v1alpha_dashboards_api.NewGetDashboardParams().WithIDOrName(name)
			resp, err := c.SemaphoreDashboardsV1alphaDashboardsAPI.GetDashboard(params)

			utils.Check(err)

			j, err := resp.Payload.MarshalYaml()

			utils.Check(err)

			fmt.Printf("apiVersion: v1alpha\n")
			fmt.Printf("kind: Dashboard\n")
			fmt.Printf("%s", j)
		}
	},
}
