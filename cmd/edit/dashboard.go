package cmd_edit

import (
	"github.com/spf13/cobra"
)

var EditDashboardCmd = &cobra.Command{
	Use:     "dashboard [name]",
	Short:   "Edit a dashboard.",
	Long:    ``,
	Aliases: []string{"dashboards dash"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		// name := args[0]

		// c := api.DefaultClient()

		// params := semaphore_dashboards_v1alpha_dashboards_api.NewGetDashboardParams().WithIDOrName(name)
		// resp, err := c.SemaphoreDashboardsV1alphaDashboardsAPI.GetDashboard(params)

		// utils.Check(err)

		// y, err := resp.Payload.MarshalBinary()

		// utils.Check(err)

		// objectName := fmt.Sprintf("Dashboard/%s", resp.Payload.Metadata.ID)
		// content := fmt.Sprintf("apiVersion: v1alpha\nkind: Dashboard\n%s", y)

		// new_content, err := handler.EditYamlInEditor(objectName, content)

		// utils.Check(err)

		// new_dash :={}
		// updated_secret, err := yaml.UnmarshalStrict(

		// utils.Check(err)

		// err = updated_secret.Update()

		// utils.Check(err)

		// fmt.Printf("Secret '%s' updated.\n", secret.Metadata.Name)
	},
}
