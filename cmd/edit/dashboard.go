package cmd_edit

import (
	"fmt"
	"strings"

	"github.com/semaphoreci/cli/api"
	"github.com/semaphoreci/cli/api/client/semaphore_dashboards_v1alpha_dashboards_api"
	"github.com/semaphoreci/cli/cmd/handler"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var EditDashboardCmd = &cobra.Command{
	Use:     "dashboard [name]",
	Short:   "Edit a dashboard.",
	Long:    ``,
	Aliases: []string{"dashboards", "dash"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := api.DefaultClient()

		params := semaphore_dashboards_v1alpha_dashboards_api.NewGetDashboardParams().WithIDOrName(name)
		resp, err := c.SemaphoreDashboardsV1alphaDashboardsAPI.GetDashboard(params)
		dashboard := resp.Payload

		utils.Check(err)

		y, err := dashboard.MarshalYaml()

		utils.Check(err)

		objectName := fmt.Sprintf("Dashboard/%s", dashboard.Metadata.ID)
		content := fmt.Sprintf("apiVersion: v1alpha\nkind: Dashboard\n%s", y)

		new_content, err := handler.EditYamlInEditor(objectName, content)
		new_content = strings.Replace(new_content, "apiVersion: v1alpha\nkind: Dashboard\n", "", -1)

		utils.Check(err)

		err = dashboard.UnmarshalYaml([]byte(new_content))

		utils.Check(err)

		fmt.Printf("[error] No upstream handler\n")
	},
}
