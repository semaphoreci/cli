package client

import (
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
	"github.com/renderedtext/sem/config"
)

type Project struct {
	Metadata struct {
		Name string
		Id   string
	}

	Spec struct {
		Repository struct {
			Url string
		}
	}
}

func InitProject(name string, repo_url string) Project {
	p := Project{}

	p.Metadata.Name = name
	p.Spec.Repository.Url = repo_url

	return p
}

func InitProjectFromYaml(data []byte) (Project, error) {
	p := Project{}

	err := yaml.Unmarshal(data, &p)

	if err != nil {
		return p, err
	}

	return p, nil
}

func (*Project) Create() {
	c := client.FromConfig()

	body, status, err := c.Post("projects", project_yaml)

	if err != nil {
		return "", errors.New(fmt.Sprintf("connecting project to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return "", errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	project_url := fmt.Sprintf(, host, name)

	return project_url, nil
}
