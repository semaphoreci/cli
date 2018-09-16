package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type ProjectV1Alpha struct {
	ApiVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Metadata   struct {
		Name string `json:"name,omitempty"`
		Id   string `json:"id,omitempty"`
	} `json:"metadata,omitempty"`

	Spec struct {
		Repository struct {
			Url string `json:"url,omitempty"`
		} `json:"repository,omitempty"`
	} `json:"spec,omitempty"`
}

func NewProjectV1Alpha(name string) ProjectV1Alpha {
	d := ProjectV1Alpha{}

	d.ApiVersion = "v1alpha"
	d.Kind = "Project"
	d.Metadata.Name = name

	return d
}

func NewProjectV1AlphaFromJson(data []byte) (*ProjectV1Alpha, error) {
	d := ProjectV1Alpha{}

	err := json.Unmarshal(data, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func NewProjectV1AlphaFromYaml(data []byte) (*ProjectV1Alpha, error) {
	d := ProjectV1Alpha{}

	err := json.Unmarshal(data, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (s *ProjectV1Alpha) ObjectName() string {
	return fmt.Sprintf("Projects/%s", s.Metadata.Name)
}

func (s *ProjectV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ProjectV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}
