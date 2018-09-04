package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd/utils"
	"github.com/renderedtext/sem/config"
	"github.com/renderedtext/sem/generators"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		RunInit(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(InitCmd)
}

func RunInit(cmd *cobra.Command, args []string) {
	repo_url, err := getGitOriginUrl()

	utils.Check(err)

	name, err := ConstructProjectName(repo_url)

	utils.Check(err)

	project := client.InitProject(name, repo_url)
	err = project.Create()

	utils.Check(err)

	if generators.PipelineFileExists() {
		err = generators.GeneratePipelineYaml()

		utils.Check(err)
	} else {
		fmt.Printf("[info] skipping .semaphore/semaphore.yml generation. It is already present in the repository.")
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
		regexp.MustCompile(`git\@github\.com:.*\/(.*).git`),
		regexp.MustCompile(`git\@github\.com:.*\/(.*)`),
		regexp.MustCompile(`git\@github\.com:.*\/(.*)`),
	}

	for _, r := range formats {
		match := r.FindStringSubmatch(repo_url)

		if len(match) >= 2 {
			return match[1], nil
		}
	}

	errTemplate := "unknown git remote format '%s'.\n"
	errTemplate += "\n"
	errTemplate += "Format must be one of the following:\n"
	errTemplate += "  - git@github.com:/<owner>/<repo_name>.git\n"
	errTemplate += "  - git@github.com:/<owner>/<repo_name>\n"

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
