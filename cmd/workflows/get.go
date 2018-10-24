package workflows

import (
	"fmt"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/cmd/utils"
)

func List(projectName string) {
	projectClient := client.NewProjectV1AlphaApi()
	project, err := projectClient.GetProject(projectName)
	utils.CheckWithMessage(err, fmt.Sprintf("project_id for project '%s' not found; '%s'", projectName, err))

	fmt.Printf("project id: %s\n", project.Metadata.Id)

	wfClient := client.NewWorkflowV1AlphaApi()
	workflows, err := wfClient.ListWorkflows(project.Metadata.Id)
	utils.Check(err)

	prettyPrint(workflows)
}

func prettyPrint(workflows *models.WorkflowListV1Alpha) {
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	fmt.Fprintln(w, "WORKFLOW ID\tINITIAL PIPELINE ID\tLABEL")

	for _, p := range workflows.Workflow {
		fmt.Fprintf(w, "%s\t%s\t%s\n", p.Id, p.InitialPplId, p.BranchName)
	}

	w.Flush()
}
