package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [KIND]",
	Short: "List resources.",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 2),
}

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

var GetSecretCmd = &cobra.Command{
	Use:     "secrets [name]",
	Short:   "Get secrets.",
	Long:    ``,
	Aliases: []string{"secret"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewSecretV1BetaApi()

		if len(args) == 0 {
			secretList, err := c.ListSecrets()

			utils.Check(err)

			const padding = 3
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

			fmt.Fprintln(w, "NAME\tAGE")

			for _, s := range secretList.Secrets {
				updateTime, err := s.Metadata.UpdateTime.Int64()

				utils.Check(err)

				fmt.Fprintf(w, "%s\t%s\n", s.Metadata.Name, utils.RelativeAgeForHumans(updateTime))
			}

			w.Flush()
		} else {
			name := args[0]

			secret, err := c.GetSecret(name)

			utils.Check(err)

			y, err := secret.ToYaml()

			utils.Check(err)

			fmt.Printf("%s", y)
		}
	},
}

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

var GetJobAllStates bool

var GetJobCmd = &cobra.Command{
	Use:     "jobs [id]",
	Short:   "Get jobs.",
	Long:    ``,
	Aliases: []string{"job"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewJobsV1AlphaApi()

		if len(args) == 0 {
			states := []string{
				"PENDING",
				"QUEUED",
				"RUNNING",
			}

			if GetJobAllStates {
				states = append(states, "FINISHED")
			}

			jobList, err := c.ListJobs(states)

			utils.Check(err)

			const padding = 3
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

			fmt.Fprintln(w, "ID\tNAME\tAGE\tSTATE\tRESULT")

			for _, j := range jobList.Jobs {
				createTime, err := j.Metadata.CreateTime.Int64()

				utils.Check(err)

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
					j.Metadata.Id,
					j.Metadata.Name,
					utils.RelativeAgeForHumans(createTime),
					j.Status.State,
					j.Status.Result)
			}

			w.Flush()
		} else {
			id := args[0]

			job, err := c.GetJob(id)

			utils.Check(err)

			y, err := job.ToYaml()

			utils.Check(err)

			fmt.Printf("%s", y)
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	getCmd.AddCommand(GetDashboardCmd)
	getCmd.AddCommand(GetSecretCmd)
	getCmd.AddCommand(GetProjectCmd)

	GetJobCmd.Flags().BoolVar(&GetJobAllStates, "all", false, "list all jobs including finished ones")
	getCmd.AddCommand(GetJobCmd)
}
