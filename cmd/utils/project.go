package utils

import (
	"fmt"

	"github.com/semaphoreci/cli/api/client"

	"log"
	"os/exec"
	"strings"
)

func GetProjectId(name string) string {
	projectClient := client.NewProjectV1AlphaApi()
	project, err := projectClient.GetProject(name)

	CheckWithMessage(err, fmt.Sprintf("project_id for project '%s' not found; '%s'", name, err))

	return project.Metadata.Id
}

func InferProjectName() (string, error) {
	originUrl, err := getGitOriginUrl()
	if err != nil {
		return "", err
	}

	log.Printf("Origin url: '%s'\n", originUrl)

	projectName, err := getProjectIdFromUrl(originUrl)
	if err != nil {
		return "", err
	}

	return projectName, nil
}

func getProjectIdFromUrl(url string) (string, error) {
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

func getGitOriginUrl() (string, error) {
	args := []string{"config", "remote.origin.url"}

	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		cmd_string := fmt.Sprintf("'%s %s'", "git", strings.Join(args, " "))
		user_msg := "You are probably not in a git directory?"
		return "", fmt.Errorf("%s failed with message: '%s'\n%s", cmd_string, err, user_msg)
	}

	return strings.TrimSpace(string(out)), nil
}
