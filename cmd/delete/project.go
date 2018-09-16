package cmd_delete

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var DeleteProjectCmd = &cobra.Command{
	Use:     "project [NAME]",
	Short:   "Delete a project.",
	Long:    ``,
	Aliases: []string{"projects", "prj"},
	Args:    cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		c := client.NewProjectV1AlphaApi()

		err := c.DeleteProject(name)

		utils.Check(err)

		fmt.Printf("Project '%s' deleted.\n", name)
	},
}
