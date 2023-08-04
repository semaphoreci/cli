package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type ProjectSecretV1 struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string `json:"kind,omitempty" yaml:"kind"`

	Metadata ProjectSecretV1Metadata `json:"metadata" yaml:"metadata"`
	Data     ProjectSecretV1Data     `json:"data" yaml:"data"`
}

type ProjectSecretV1EnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type ProjectSecretV1File struct {
	Path    string `json:"path" yaml:"path"`
	Content string `json:"content" yaml:"content"`
}

type ProjectSecretV1Data struct {
	EnvVars []ProjectSecretV1EnvVar `json:"env_vars" yaml:"env_vars"`
	Files   []ProjectSecretV1File   `json:"files" yaml:"files"`
}

type ProjectSecretV1Metadata struct {
	Name            string      `json:"name,omitempty" yaml:"name,omitempty"`
	Id              string      `json:"id,omitempty" yaml:"id,omitempty"`
	CreateTime      json.Number `json:"create_time,omitempty,string" yaml:"create_time,omitempty"`
	UpdateTime      json.Number `json:"update_time,omitempty,string" yaml:"update_time,omitempty"`
	ProjectIdOrName string      `json:"project_id_or_name,omitempty" yaml:"project_id_or_name,omitempty"`
}

func NewProjectSecretV1(name string, envVars []ProjectSecretV1EnvVar, files []ProjectSecretV1File) ProjectSecretV1 {
	s := ProjectSecretV1{}

	s.setApiVersionAndKind()
	s.Metadata.Name = name
	s.Data.EnvVars = envVars
	s.Data.Files = files

	return s
}

func NewProjectSecretV1FromJson(data []byte) (*ProjectSecretV1, error) {
	s := ProjectSecretV1{}

	err := json.Unmarshal(data, &s)

	if err != nil {
		return nil, err
	}

	s.setApiVersionAndKind()

	return &s, nil
}

func NewProjectSecretV1FromYaml(data []byte) (*ProjectSecretV1, error) {
	s := ProjectSecretV1{}

	err := yaml.UnmarshalStrict(data, &s)

	if err != nil {
		return nil, err
	}

	s.setApiVersionAndKind()

	return &s, nil
}

func (s *ProjectSecretV1) setApiVersionAndKind() {
	s.ApiVersion = "v1"
	s.Kind = "ProjectSecret"
}

func (s *ProjectSecretV1) ObjectName() string {
	return fmt.Sprintf("Project/%s/Secrets/%s", s.Metadata.ProjectIdOrName, s.Metadata.Name)
}

func (s *ProjectSecretV1) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *ProjectSecretV1) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}
