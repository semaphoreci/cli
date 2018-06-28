package cmd

import (
	"fmt"
  "io/ioutil"

	"github.com/spf13/cobra"
	"github.com/ghodss/yaml"
	"github.com/renderedtext/sem/client"
	"github.com/spf13/viper"
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

  name := "test-125"
  url := "git@github.com:shiroyasha/test.git"

  project, err := yaml.YAMLToJSON([]byte(fmt.Sprintf(project_template, name, url)))

  if err != nil {
    fmt.Printf("%v", err)
  }

	body, status, _ := c.Post("projects", project)

  if status != 200 {
    fmt.Println(string(body))
    return
  }

  err = ioutil.WriteFile(".semaphore.yml", []byte(semaphore_yaml_template), 0644)

  if err != nil {
    fmt.Printf("%v", err)
    return
  }

  project_url := fmt.Sprintf("https://%s/projects/%s", viper.GetString("host"), name)

  fmt.Printf("Project \"%s\" initialized on Semaphore. %s\n", name, project_url)
  fmt.Println("")
  fmt.Println("Execute: git add .semaphore.yml && git commit -m \"First pipeline\" && git push")
}
