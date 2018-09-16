package cmd_get

import (
	"fmt"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var GetDashboardCmd = &cobra.Command{
	Use:     "dashboards [name]",
	Short:   "Get dashboards.",
	Long:    ``,
	Aliases: []string{"dashboard", "dash"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewDashboardV1AlphaApi()

		if len(args) == 0 {
			dashList, err := c.ListDashboards()

			utils.Check(err)

			const padding = 3
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

			fmt.Fprintln(w, "NAME\tAGE")

			for _, d := range dashList.Dashboards {
				updateTime, err := d.Metadata.UpdateTime.Int64()

				utils.Check(err)

				fmt.Fprintf(w, "%s\t%s\n", d.Metadata.Name, utils.RelativeAgeForHumans(updateTime))
			}

			w.Flush()
		} else {
			name := args[0]

			dash, err := c.GetDashboard(name)

			utils.Check(err)

			y, err := dash.ToYaml()

			utils.Check(err)

			fmt.Printf("%s", y)
		}
	},
}
