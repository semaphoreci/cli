package cmd

import (
	"fmt"
  "io/ioutil"
  "os"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/ghodss/yaml"
	"github.com/renderedtext/sem/client"
	"github.com/spf13/viper"
	"github.com/tcnksm/go-gitconfig"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project",
	Long: ``,

	Run: func(cmd *cobra.Command, args []string) {
    RunInit(cmd, args)
	},
}

var project_template = `
apiVersion: v1alpha
kind: Project
metadata:
  name: %s
spec:
  repository:
    url: "%s"
`

var semaphore_yaml_template = `
version: "v3.0"
name: My first pipeline
semaphore_image: standard
blocks:
  - name: "Stage 1"
    build:
      prologue:
        commands:
          - checkout && cd $SEMAPHORE_GIT_DIR
      epilogue:
        commands:
          - echo "Yay job finished"
          - echo $SEMAPHORE_JOB_RESULT
      env_vars:
        - name: VAR1
          value: Environment Variable 1
        - name: PI
          value: "3.14159"
      jobs:
      - name: Just ls
        commands:
          - pwd
          - echo "test"
          - ls /etc

      - name: List files
        commands:
          - echo "First env var -> $VAR1"
          - echo "My files:"
          - ls -lah

  - name: "Stage 2"
    build:
      jobs:
      - name: Echo job
        commands:
          - checkout
          - cd $SEMAPHORE_GIT_DIR
          - pwd
          - echo $SEMAPHORE_PIPELINE_ID
          - echo "Hello from $SEMAPHORE_JOB_ID"
`

func init() {
	rootCmd.AddCommand(initCmd)
}

func RunInit(cmd *cobra.Command, args []string) {
	c := client.FromConfig()

  repo_url, err := gitconfig.OriginURL()

  check(err, "Failed to extract git origin from gitconfig")

  re := regexp.MustCompile(`git\@github\.com:.*\/(.*).git`)
  match := re.FindStringSubmatch(repo_url)

	name := match[1]
	host := viper.GetString("host")
  project_url := fmt.Sprintf("https://%s/projects/%s", host, name)

  check(err, "Failed to construct project name")

  err = ioutil.WriteFile(".semaphore.yml", []byte(semaphore_yaml_template), 0644)

  check(err, "Failed to create .semaphore.yml")

  project, err := yaml.YAMLToJSON([]byte(fmt.Sprintf(project_template, name, repo_url)))

  check(err, "Failed to connect project to Semaphore")

	body, status, err := c.Post("projects", project)

  check(err, "Failed to connect project to Semaphore")

  if status != 200 {
		fmt.Fprintf(os.Stderr, "%s\n", body)
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

    os.Exit(1)
  }

  fmt.Printf("Project is created. You can find it at %s.\n", project_url)
  fmt.Println("")
  fmt.Printf("To run our first pipeline execute:")
  fmt.Printf("git add .semaphore.yml && git commit -m \"First pipeline\" && git push")
}
