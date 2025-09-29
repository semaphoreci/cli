package utils

import (
	"fmt"

	"log"
	"os/exec"
	"strings"

	"github.com/semaphoreci/cli/api/client"
)

func GetProjectId(name string) string {
	projectClient := client.NewProjectV1AlphaApi()
	project, err := projectClient.GetProject(name)

	CheckWithMessage(err, fmt.Sprintf("project_id for project '%s' not found; '%s'", name, err))

	return project.Metadata.Id
}

func InferProjectName() (string, error) {
	originURLs, err := getAllGitRemoteURLs()
	if err != nil {
		return "", err
	}

	var projectName string

	for _, originURL := range originURLs {
		log.Printf("Origin url: '%s'\n", originURL)

		projectName, err = getProjectIdFromUrl(originURL)
		if err != nil {
			log.Printf("no project found for remote %s", originURL)
		}
		if projectName != "" {
			return projectName, nil
		}
	}

	return "", fmt.Errorf("Unable to find project for any configured remotes")
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
