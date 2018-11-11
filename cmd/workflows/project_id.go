package workflows

import (
	"fmt"

	"github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"

	"os/exec"
	"strings"
)

func GetProjectId(name string) string {
	projectClient := client.NewProjectV1AlphaApi()
	project, err := projectClient.GetProject(name)

	utils.CheckWithMessage(err, fmt.Sprintf("project_id for project '%s' not found; '%s'", name, err))

	return project.Metadata.Id
}

func GetProjectIdFromUrl(url string) (string, error) {
	projectClient := client.NewProjectV1AlphaApi()
	projects, err := projectClient.ListProjects()

	if err != nil {
		return "", fmt.Errorf("getting project list failed '%s'", err)
	}

	projectName := ""
	for _, p := range projects.Projects {
		if p.Spec.Repository.Url == url {
			projectName = p.Metadata.Name
			break
		}
	}

	if projectName == "" {
		return "", fmt.Errorf("project with url '%s' not found in this org", url)
	}


	return projectName, nil
}

func GetGitOriginUrl() (string, error) {
	args := []string{"config", "remote.origin.url"}

	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("'git config remote.origin.url' failed with message: '%s'", err)
	}

	return strings.TrimSpace(string(out)), nil
}
