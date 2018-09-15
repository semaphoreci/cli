package cmd_delete

import (
	"github.com/spf13/cobra"
)

var DeleteDashboardCmd = &cobra.Command{
	Use:     "dashboard [NAME]",
	Short:   "Delete a dashboard.",
	Long:    ``,
	Aliases: []string{"dashboard", "dash"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		// name := args[0]

		// c := api.DefaultClient()

		// params := semaphore_dashboards_v1alpha_dashboards_api.NewDeleteDashboardParams()
		// params.SetIDOrName(name)

		// _, err := c.SemaphoreDashboardsV1alphaDashboardsAPI.DeleteDashboard(params)

		// if err != nil {
		// 	utils.Check(err)
		// }

		// fmt.Printf("Dashboard '%s' deleted.\n", name)
	},
}
