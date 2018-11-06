package workflows

import (
	"fmt"

	"github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
)

func GetProjectId(name string) string {
	projectClient := client.NewProjectV1AlphaApi()
	project, err := projectClient.GetProject(name)

	utils.CheckWithMessage(err, fmt.Sprintf("project_id for project '%s' not found; '%s'", name, err))

	return project.Metadata.Id
}

func GetProjectIdFromUrl(url string) string {
	projectClient := client.NewProjectV1AlphaApi()
	project, err := projectClient.ListProjects()

	utils.CheckWithMessage(err, fmt.Sprintf("project_id for project  not found; '%s'",  err))

	fmt.Print(project)
	return "as"
}
