package cmd

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/api/uuid"
	"github.com/semaphoreci/cli/cmd/deployment_targets"
	"github.com/semaphoreci/cli/cmd/pipelines"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/semaphoreci/cli/cmd/workflows"
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

			if err := w.Flush(); err != nil {
				fmt.Printf("Error flushing dashboards: %v\n", err)
			}
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
		projectID := GetProjectID(cmd)

		if projectID == "" {
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

				if err := w.Flush(); err != nil {
					fmt.Printf("Error flushing secrets: %v\n", err)
				}
			} else {
				name := args[0]

				secret, err := c.GetSecret(name)

				utils.Check(err)

				y, err := secret.ToYaml()

				utils.Check(err)

				fmt.Printf("%s", y)
			}
		} else {
			c := client.NewProjectSecretV1Api(projectID)

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

				if err := w.Flush(); err != nil {
					fmt.Printf("Error flushing secrets: %v\n", err)
				}
			} else {
				name := args[0]

				secret, err := c.GetSecret(name)

				utils.Check(err)

				y, err := secret.ToYaml()

				utils.Check(err)

				fmt.Printf("%s", y)
			}
		}
	},
}

var GetAgentTypeCmd = &cobra.Command{
	Use:     "agent_types [name]",
	Short:   "Get self-hosted agent types.",
	Long:    ``,
	Aliases: []string{"agent_type", "agenttype", "agenttypes", "agentTypes", "agentType"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewAgentTypeApiV1AlphaApi()

		if len(args) == 0 {
			agentTypeList, err := c.ListAgentTypes()

			utils.Check(err)

			const padding = 3
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

			fmt.Fprintln(w, "NAME\tAGE")

			for _, a := range agentTypeList.AgentTypes {
				updateTime, err := a.Metadata.UpdateTime.Int64()
				utils.Check(err)
				fmt.Fprintf(w, "%s\t%s\n", a.Metadata.Name, utils.RelativeAgeForHumans(updateTime))
			}

			if err := w.Flush(); err != nil {
				fmt.Printf("Error flushing agent types: %v\n", err)
			}
		} else {
			name := args[0]

			secret, err := c.GetAgentType(name)

			utils.Check(err)

			y, err := secret.ToYaml()

			utils.Check(err)

			fmt.Printf("%s", y)
		}
	},
}

var GetAgentsCmd = &cobra.Command{
	Use:     "agents",
	Short:   "Get self-hosted agents.",
	Long:    ``,
	Aliases: []string{"agent"},
	Args:    cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewAgentApiV1AlphaApi()

		if len(args) == 0 {
			agentType, err := cmd.Flags().GetString("agent-type")
			utils.Check(err)

			agents, err := getAllAgents(c, agentType)
			utils.Check(err)

			if len(agents) == 0 {
				fmt.Fprintln(os.Stdout, "No agents found")
				return
			}

			const padding = 3
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

			fmt.Fprintln(w, "NAME\tTYPE\tSTATE\tAGE")
			for _, a := range agents {
				connectedAt, err := a.Metadata.ConnectedAt.Int64()
				utils.Check(err)
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
					a.Metadata.Name,
					a.Metadata.Type,
					a.Status.State,
					utils.RelativeAgeForHumans(connectedAt),
				)
			}

			if err := w.Flush(); err != nil {
				fmt.Printf("Error flushing agents: %v\n", err)
			}
		} else {
			name := args[0]

			agent, err := c.GetAgent(name)

			utils.Check(err)

			y, err := agent.ToYaml()

			utils.Check(err)

			fmt.Printf("%s", y)
		}
	},
}

func getAllAgents(client client.AgentApiV1AlphaApi, agentType string) ([]models.AgentV1Alpha, error) {
	agents := []models.AgentV1Alpha{}
	cursor := ""

	for {
		agentList, err := client.ListAgents(agentType, cursor)
		if err != nil {
			return nil, err
		}

		agents = append(agents, agentList.Agents...)
		if agentList.Cursor == "" {
			break
		}

		cursor = agentList.Cursor
	}

	return agents, nil
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

			if err := w.Flush(); err != nil {
				fmt.Printf("Error flushing projects: %v\n", err)
			}
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

			if err := w.Flush(); err != nil {
				fmt.Printf("Error flushing jobs: %v\n", err)
			}
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

var GetPplFollow bool

var GetPplCmd = &cobra.Command{
	Use:     "pipelines [id]",
	Short:   "Get pipelines.",
	Long:    ``,
	Aliases: []string{"pipeline", "ppl"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			projectID := getPrj(cmd)

			pipelines.List(projectID)
		} else {
			id := args[0]
			pipelines.Describe(id, GetPplFollow)
		}
	},
}

var GetWfCmd = &cobra.Command{
	Use:     "workflows [id]",
	Short:   "Get workflows.",
	Long:    ``,
	Aliases: []string{"workflow", "wf"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		projectID := getPrj(cmd)

		if len(args) == 0 {
			workflows.List(projectID)
		} else {
			wfID := args[0]
			workflows.Describe(projectID, wfID)
		}
	},
}

var GetDTCmd = &cobra.Command{
	Use:     "deployment_targets [id or name]",
	Short:   "Get deployment targets.",
	Long:    ``,
	Aliases: models.DeploymentTargetCmdAliases,
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		getHistory, err := cmd.Flags().GetBool("history")
		utils.Check(err)
		targetName, err := cmd.Flags().GetString("name")
		utils.Check(err)
		targetId, err := cmd.Flags().GetString("id")
		utils.Check(err)
		if len(args) == 1 {
			if uuid.IsValid(args[0]) {
				targetId = args[0]
			} else {
				targetName = args[0]
			}
		}
		if getHistory {
			if targetId != "" {
				deployment_targets.History(targetId, cmd)
			} else if targetName != "" {
				deployment_targets.HistoryByName(targetName, getPrj(cmd), cmd)
			}
		} else {
			if targetId != "" {
				deployment_targets.Describe(targetId)
			} else if targetName != "" {
				deployment_targets.DescribeByName(targetName, getPrj(cmd))
			} else {
				deployment_targets.List(getPrj(cmd))
			}
		}
	},
}

func GetProjectID(cmd *cobra.Command) string {
	projectID, err := cmd.Flags().GetString("project-id")
	if projectID != "" {
		return projectID
	}

	projectName, err := cmd.Flags().GetString("project-name")
	utils.Check(err)

	if projectName == "" {
		return ""
	}

	projectID = utils.GetProjectId(projectName)
	log.Printf("Project ID: %s\n", projectID)

	return projectID
}

func getPrj(cmd *cobra.Command) string {
	projectID, err := cmd.Flags().GetString("project-id")
	if projectID != "" {
		return projectID
	}

	projectName, err := cmd.Flags().GetString("project-name")
	utils.Check(err)

	if projectName == "" {
		projectName, err = utils.InferProjectName()
		utils.Check(err)
	}
	log.Printf("Project Name: %s\n", projectName)

	projectID = utils.GetProjectId(projectName)
	log.Printf("Project ID: %s\n", projectID)

	return projectID
}

func init() {
	RootCmd.AddCommand(getCmd)

	getNotificationCmd := NewGetNotificationCmd()

	getCmd.AddCommand(GetDashboardCmd)
	getCmd.AddCommand(getNotificationCmd)
	getCmd.AddCommand(GetProjectCmd)
	getCmd.AddCommand(GetAgentTypeCmd)

	GetAgentsCmd.Flags().StringP("agent-type", "t", "",
		"agent type; if specified, returns only agents for this agent type")
	getCmd.AddCommand(GetAgentsCmd)

	GetSecretCmd.Flags().StringP("project-name", "p", "",
		"project name; if specified will get secret from project level, otherwise organization secret")
	GetSecretCmd.Flags().StringP("project-id", "i", "",
		"project id; if specified will get secret from project level, otherwise organization secret")
	getCmd.AddCommand(GetSecretCmd)

	GetJobCmd.Flags().BoolVar(&GetJobAllStates, "all", false, "list all jobs including finished ones")
	getCmd.AddCommand(GetJobCmd)

	GetPplCmd.Flags().BoolVarP(&GetPplFollow, "follow", "f", false,
		"repeat get until pipeline reaches terminal state")
	GetPplCmd.Flags().StringP("project-name", "p", "",
		"project name; if not specified will be inferred from git origin")
	GetPplCmd.Flags().StringP("project-id", "i", "",
		"project id; if not specified will be inferred from git origin")
	getCmd.AddCommand(GetPplCmd)

	getCmd.AddCommand(GetWfCmd)
	GetWfCmd.Flags().StringP("project-name", "p", "",
		"project name; if not specified will be inferred from git origin")
	GetWfCmd.Flags().StringP("project-id", "i", "",
		"project id; if not specified will be inferred from git origin")

	getCmd.AddCommand(GetDTCmd)
	GetDTCmd.Flags().StringP("project-name", "p", "",
		"project name; if not specified will be inferred from git origin")
	GetDTCmd.Flags().StringP("project-id", "i", "",
		"project id; if not specified will be inferred from git origin")
	GetDTCmd.Flags().StringP("id", "t", "", "target id")
	GetDTCmd.Flags().StringP("name", "n", "", "target name")
	GetDTCmd.Flags().BoolP("history", "s", false, "get deployment target history")
	GetDTCmd.Flags().Lookup("history").NoOptDefVal = "true"
	GetDTCmd.Flags().StringP("after", "a", "", "show deployment history after the timestamp")
	GetDTCmd.Flags().StringP("before", "b", "", "show deployment history before the timestamp")
	GetDTCmd.Flags().StringP("git-ref-type", "g", "", "git reference type: branch, tag, pr")
	GetDTCmd.Flags().StringP("git-ref-label", "l", "", "git reference label: branch or tag name")
	GetDTCmd.Flags().StringArrayP("parameter", "q", []string{}, "show deployment history of deployment targets with provided bookmark parameters")
	GetDTCmd.Flags().StringP("triggered-by", "u", "", "show deployment history triggered by specific user or promotion")
}
