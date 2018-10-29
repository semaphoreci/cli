package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/semaphoreci/cli/config"
	"github.com/semaphoreci/cli/generators"
	"github.com/spf13/cobra"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	gitconfig "github.com/tcnksm/go-gitconfig"
)

var flagProjectName string
var flagRepoUrl string

func InitCmd() cobra.Command {
	cmd := cobra.Command{
		Use:   "init",
		Short: "Initialize a project",
		Long:  ``,

		Run: func(cmd *cobra.Command, args []string) {
			RunInit(cmd, args)
		},
	}

	cmd.Flags().StringVar(&flagRepoUrl, "repo-url", "", "explicitly set the repository url, if not set it is extracted from local git repository")
	cmd.Flags().StringVar(&flagProjectName, "project-name", "", "explicitly set the project name, if not set it is extracted from the repo-url")

	return cmd
}

func init() {
	cmd := InitCmd()

	RootCmd.AddCommand(&cmd)
}

func RunInit(cmd *cobra.Command, args []string) {
	var err error
	var name string
	var repoUrl string

	if flagRepoUrl != "" {
		repoUrl = flagRepoUrl
	} else {
		repoUrl, err = getGitOriginUrl()

		utils.Check(err)
	}

	if flagProjectName != "" {
		name = flagProjectName
	} else {
		name, err = ConstructProjectName(repoUrl)

		utils.Check(err)
	}

	c := client.NewProjectV1AlphaApi()
	projectModel := models.NewProjectV1Alpha(name)
	projectModel.Spec.Repository.Url = repoUrl

	project, err := c.CreateProject(&projectModel)

	utils.Check(err)

	if generators.PipelineFileExists() {
		fmt.Printf("[info] skipping .semaphore/semaphore.yml generation. It is already present in the repository.\n\n")
	} else {
		err = generators.GeneratePipelineYaml()

		utils.Check(err)
	}

	fmt.Printf("Project is created. You can find it at https://%s/projects/%s.\n", config.GetHost(), project.Metadata.Name)
	fmt.Println("")
	fmt.Printf("To run your first pipeline execute:\n")
	fmt.Println("")
	fmt.Printf("  git add .semaphore/semaphore.yml && git commit -m \"First pipeline\" && git push\n")
	fmt.Println("")
}

func ConstructProjectName(repo_url string) (string, error) {
	formats := []*regexp.Regexp{
		regexp.MustCompile(`git\@github\.com:/.*\/(.*).git`),
		regexp.MustCompile(`git\@github\.com:.*\/(.*).git`),
		regexp.MustCompile(`git\@github\.com:/.*\/(.*)`),
		regexp.MustCompile(`git\@github\.com:.*\/(.*)`),
		regexp.MustCompile(`https://github.com/.*\/(.*).git`),
		regexp.MustCompile(`https://github.com/.*\/(.*)`),
		regexp.MustCompile(`http://github.com/.*\/(.*).git`),
		regexp.MustCompile(`http://github.com/.*\/(.*)`),
	}

	for _, r := range formats {
		match := r.FindStringSubmatch(repo_url)

		if len(match) >= 2 {
			return match[1], nil
		}
	}

	errTemplate := "unsupported git remote format '%s'.\n"
	errTemplate += "\n"
	errTemplate += "Format must be one of the following:\n"
	errTemplate += "  - git@github.com:<owner>/<repo_name>.git\n"
	errTemplate += "  - git@github.com:<owner>/<repo_name>\n"
	errTemplate += "  - https://github.com/<owner>/<repo_name>\n"
	errTemplate += "  - https://github.com/<owner>/<repo_name>.git\n"
	errTemplate += "\n"
	errTemplate += "To add a project with an alternative git url, use the --repo-url flag:\n"
	errTemplate += "  - sem init --repo-url git@github.com:<owner>/<repo_name>.git\n"

	return "", errors.New(fmt.Sprintf(errTemplate, repo_url))
}

func getGitOriginUrl() (string, error) {
	if flag.Lookup("test.v") == nil {
		if _, err := os.Stat(".git"); os.IsNotExist(err) {
			return "", errors.New("not a git repository")
		}

		return gitconfig.OriginURL()
	} else {
		return "git@github.com:/renderedtext/something.git", nil
	}
}
