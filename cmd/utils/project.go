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

type GitRemote struct {
	Name string
	URL string
	Project models.ProjectV1Alpha
}

type GitRemoteList []GitRemote

func (grl GitRemoteList) Contains(remoteNameOrUrl string) bool {
	for _, gitRemote := range grl {
		if gitRemote.Name == remoteNameOrUrl || gitRemote.URL == remoteNameOrUrl {
			return true
		}
	}
	return false
}

func (grl GitRemoteList) Get(remoteNameOrUrl string) (*GitRemote, error) {
	for _, gitRemote := range grl {
		if gitRemote.Name == remoteNameOrUrl || gitRemote.URL == remoteNameOrUrl {
			return &gitRemote, nil
		}
	}
	return &GitRemote{}, fmt.Errorf("no remote matching %s found in remote list")
}

func (grl GitRemoteList) URLs() []string {
	urls := []string{}
	for _, gitRemote := range grl {
		urls = append(urls, gitRemote.URL)
	}
	return urls
}

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
	// Note that getAllGitRemotesAndProjects will only return remotes
	// where the URL of that remote is configured in a project we got
	// from the API, so this list is a list of remotes valid projects
	// configured. All we have to do now is pick one.
	gitRemotes, err := getAllGitRemotesAndProjects()
	if err != nil {
		return models.ProjectV1Alpha{}, err
	}

	// If the user is using GitHub and has run `gh repo set-default`, then
	// we can get that 'base' remote name and see if we have a project for
	// it.
	ghBaseRemoteName, err := getGitHubBaseRemoteName()
	if err != nil {
		log.Printf("tried looking for a `gh` base repo configuration, but found none.")
	} else {
		gitRemote, err := gitRemotes.Get(ghBaseRemoteName)
		if err == nil {
			return gitRemote.Project, nil
		}
	}

	// If we only got one remote with a configured project, return it;
	// alternately, if we got multiple, return the alphabetically first
	// one.
	if len(gitRemotes) == 1 {
		return gitRemotes[0].Project, nil
	}

	// If we got an "origin" remote or an "upstream" remote, return that (in
	// that order of preference)
	for _, remoteName := range []string{"origin", "upstream"} {
		remote, err := gitRemotes.Get(remoteName)
		if err == nil {
			return remote.Project, nil
		}
	}

	// At this point, we have multiple remotes configured, all of which have
	// a project configured in Semaphore, none of which are named "origin" or
	// "upstream", and none of which are set as the gh base repo. The *most likely*
	// explanation here is that the user has the same repo URL configured multiple
	// times, or they're doing something extremely unusual. I'm not sure we can
	// make the correct decision here.
	allUrls := []string{}
	for _, url := range gitRemotes.URLs() {
		if !slices.Contains(allUrls, url) {
			allUrls = append(allUrls, url)
		}
	}

	// Okay, there's only one URL so we can just pick the relevant project.
	if len(allUrls) == 1 {
		return gitRemotes[0].Project, nil
	}

	// At this point we'd just be guessing, so let's give up
	return models.ProjectV1Alpha{}, fmt.Errorf("found %d remotes with %d different URLs but cannot determine the correct one", len(gitRemotes), len(allUrls))
}

// getGitHubBaseRemoteName checks to see if the `gh` cli tool has set a default
// remote for this repository. If not, or if we're not using Github at all, we
// can just ignore the error.
func getGitHubBaseRemoteName() (string, error) {
	args := []string{"config", "--local", "--get-regexp", "gh-resolved"}
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		cmd_string := fmt.Sprintf("'%s %s'", "git", strings.Join(args, " "))
		user_msg := "You are probably not in a git directory?"
		return "", fmt.Errorf("%s failed with message: '%s'\n%s", cmd_string, err, user_msg)
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	if len(lines) == 0 {
		return "", fmt.Errorf("no GitHub base remote configured for this repository")
	}
	if len(lines) > 1 {
		return "", fmt.Errorf("got multiple lines when looking for GitHub base remote")
	}

	fields := strings.Fields(lines[0])
	remoteName := strings.Split(fields[0], ".")[1]
	return remoteName, nil
}

func getAllGitRemotesAndProjects() (GitRemoteList, error) {
	args := []string{"config", "--local", "--get-regexp", "remote\\..*\\.url"}
	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		cmd_string := fmt.Sprintf("'%s %s'", "git", strings.Join(args, " "))
		user_msg := "You are probably not in a git directory?"
		return GitRemoteList{}, fmt.Errorf("%s failed with message: '%s'\n%s", cmd_string, err, user_msg)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")

	projectClient := client.NewProjectV1AlphaApi()
	projects, err := projectClient.ListProjects()

	remotes := GitRemoteList{}

	for _, line := range lines {
		fields := strings.Fields(line)
		keyFields := strings.Split(line, ".")
		remoteName := keyFields[1]
		url := fields[1]

		for _, proj := range projects.Projects {
			if proj.Spec.Repository.Url == url {
				remotes = append(remotes, GitRemote{Name: remoteName, URL: url, Project: proj})
				break
			}
		}

	}

	return remotes, nil
}
