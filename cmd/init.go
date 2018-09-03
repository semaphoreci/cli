package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd/utils"
	"github.com/renderedtext/sem/config"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		RunInit(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

const semaphore_yaml_template = `version: v1.0
name: First pipeline example
agent:
  machine:
    type: e1-standard-2
    os_image: ubuntu1804

blocks:
  - name: "Build"
    task:
      env_vars:
        - name: APP_ENV
          value: prod
      jobs:
      - name: Docker build
        commands:
          - checkout
          - ls -1
          - echo $APP_ENV
          - echo "Docker build..."
          - echo "done"

  - name: "Smoke tests"
    task:
      jobs:
      - name: Smoke
        commands:
          - checkout
          - echo "make smoke"

  - name: "Unit tests"
    task:
      jobs:
      - name: RSpec
        commands:
          - checkout
          - echo "make rspec"

      - name: Lint code
        commands:
          - checkout
          - echo "make lint"

      - name: Check security
        commands:
          - checkout
          - echo "make security"

  - name: "Integration tests"
    task:
      jobs:
      - name: Cucumber
        commands:
          - checkout
          - echo "make cucumber"

  - name: "Push Image"
    task:
      jobs:
      - name: Push
        commands:
          - checkout
          - echo "make docker.push"
`

func RunInit(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		utils.Fail("not a git repository")
	}

	if _, err := os.Stat(".semaphore/semaphore.yml"); err == nil {
		utils.Fail(".semaphore/semaphore.yml is already present in the repository")
	}

	repo_url, err := gitconfig.OriginURL()

	utils.CheckWithMessage(err, "failed to extract remote from git configuration")

	name, err := ConstructProjectName(repo_url)

	utils.Check(err)

	host := config.GetHost()

	project_url := registerProjectOnSemaphore(name, host, repo_url)

	createInitialSemaphoreYaml()

	fmt.Printf("Project is created. You can find it at %s.\n", project_url)
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

func createInitialSemaphoreYaml() {
	if _, err := os.Stat(".semaphore"); err != nil {
		err := os.Mkdir(".semaphore", 0755)

		utils.Check(err)
	}

	err := ioutil.WriteFile(".semaphore/semaphore.yml", []byte(semaphore_yaml_template), 0644)

	utils.CheckWithMessage(err, "failed to create .semaphore/semaphore.yml")
}

const project_template = `
apiVersion: v1alpha
kind: Project
metadata:
  name: %s
spec:
  repository:
    url: "%s"
`

func registerProjectOnSemaphore(name string, host string, repo_url string) string {
	c := client.FromConfig()

	project, err := yaml.YAMLToJSON([]byte(fmt.Sprintf(project_template, name, repo_url)))

	utils.CheckWithMessage(err, "connecting project to Semaphore failed")

	body, status, err := c.Post("projects", project)

	utils.CheckWithMessage(err, "connecting project to Semaphore failed")

	if status != 200 {
		utils.Fail(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return fmt.Sprintf("https://%s/projects/%s", host, name)
}
