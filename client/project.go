package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
)

type Project struct {
	Metadata struct {
		Name string `json:"name,omitempty"`
		Id   string `json:"id,omitempty"`
	} `json:"metadata,omitempty"`

	Spec struct {
		Repository struct {
			Url string `json:"url,omitempty"`
		} `json:"repository,omitempty"`
	} `json:"spec,omitempty"`
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

func (p *Project) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Project) Create() error {
	c := FromConfig()

	json_body, err := p.ToJson()

	if err != nil {
		return errors.New(fmt.Sprintf("failed to serialize project object '%s'", err))
	}

	body, status, err := c.Post("projects", json_body)

	if err != nil {
		return errors.New(fmt.Sprintf("creating project to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return nil
}
