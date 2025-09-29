package utils

import (
	"fmt"
	"slices"

	"log"
	"os/exec"
	"strings"

	"github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/api/models"
)

func GetProjectId(name string) string {
	projectClient := client.NewProjectV1AlphaApi()
	project, err := projectClient.GetProject(name)

	CheckWithMessage(err, fmt.Sprintf("project_id for project '%s' not found; '%s'", name, err))

	return project.Metadata.Id
}

func InferProjectName() (string, error) {
	project, err := InferProject()
	if err != nil {
		return "", err
	}
	return project.Metadata.Name, nil
}

func InferProject() (models.ProjectV1Alpha, error) {
	originURLs, err := getAllGitRemoteURLs()
	if err != nil {
		return models.ProjectV1Alpha{}, err
	}

	project, err := getProjectFromUrls(originURLs)
	if err != nil {
		log.Printf("no project found for any configured remotes (%s)", strings.Join(originURLs, ", "))
	}

	return project, nil
}

func getProjectIdFromUrl(url string) (string, error) {
	project, err := getProjectFromUrls([]string{url})
	if err != nil {
		return "", fmt.Errorf("project with url %s not found in this org", url)
	}
	return project.Metadata.Id, nil

}

func getProjectFromUrls(urls []string) (models.ProjectV1Alpha, error) {
	projectClient := client.NewProjectV1AlphaApi()
	projects, err := projectClient.ListProjects()

	if err != nil {
		return models.ProjectV1Alpha{}, fmt.Errorf("getting project list failed '%s'", err)
	}

	for _, p := range projects.Projects {
		if slices.Contains(urls, p.Spec.Repository.Url) {
			return p, nil
		}
	}

	return models.ProjectV1Alpha{}, fmt.Errorf("project with urls '%s' not found in this org", strings.Join(urls, ", "))
}

func getGitRemotes() ([]string, error) {
	args := []string{"remote"}
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		cmd_string := fmt.Sprintf("'%s %s'", "git", strings.Join(args, " "))
		user_msg := "You are probably not in a git directory?"
		return []string{}, fmt.Errorf("%s failed with message: '%s'\n%s", cmd_string, err, user_msg)
	}
	return strings.Split(strings.TrimSpace(string(out)), "\n"), nil
}

func getGitRemoteUrl(remote string) (string, error) {
	args := []string{"config", fmt.Sprintf("remote.%s.url", remote)}

	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		cmd_string := fmt.Sprintf("'%s %s'", "git", strings.Join(args, " "))
		user_msg := "You are probably not in a git directory?"
		return "", fmt.Errorf("%s failed with message: '%s'\n%s", cmd_string, err, user_msg)
	}

	return strings.TrimSpace(string(out)), nil
}

func getAllGitRemoteURLs() ([]string, error) {
	gitRemotes, err := getGitRemotes()
	if err != nil {
		return gitRemotes, err
	}
	var gitRemoteURLs []string
	for _, remote := range gitRemotes {
		remoteURL, err := getGitRemoteUrl(remote)
		if err != nil {
			log.Printf("could not get URL for remote %s", remote)
		} else {
			gitRemoteURLs = append(gitRemoteURLs, remoteURL)
		}
	}
	return gitRemoteURLs, nil
}

func getGitOriginUrl() (string, error) {
	return getGitRemoteUrl("origin")
}
