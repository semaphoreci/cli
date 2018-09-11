package generators

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

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

func GeneratePipelineYaml() error {
	path := ".semaphore/semaphore.yml"

	if _, err := os.Stat(".semaphore"); err != nil {
		err := os.Mkdir(".semaphore", 0755)

		if err != nil {
			return errors.New(fmt.Sprintf("failed to create .semaphore directory '%s'", err))
		}
	}

	err := ioutil.WriteFile(".semaphore/semaphore.yml", []byte(semaphore_yaml_template), 0644)

	if err == nil {
		return nil
	} else {
		return errors.New(fmt.Sprintf("failed to create %s file '%s'", path, err))
	}
}

func PipelineFileExists() bool {
	_, err := os.Stat(".semaphore/semaphore.yml")

	return err == nil
}
