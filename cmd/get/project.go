package cmd_get

import (
	"fmt"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var GetProjectCmd = &cobra.Command{
	Use:     "projects [name]",
	Short:   "Get projects.",
	Long:    ``,
	Aliases: []string{"project", "prj"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewProjectV1AlphaApi()

		if len(args) == 0 {
			projectList, err := c.ListProjects()

			utils.Check(err)

			const padding = 3
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

			fmt.Fprintln(w, "NAME\tREPOSITORY")

			for _, p := range projectList.Projects {
				fmt.Fprintf(w, "%s\t%s\n", p.Metadata.Name, p.Spec.Repository.Url)
			}

			w.Flush()
		} else {
			name := args[0]

			project, err := c.GetProject(name)

			utils.Check(err)

			y, err := project.ToYaml()

			utils.Check(err)

			fmt.Printf("%s", y)
		}
	},
}
