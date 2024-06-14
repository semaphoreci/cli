package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type SecretV1Beta struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string `json:"kind,omitempty" yaml:"kind"`

	Metadata  SecretV1BetaMetadata   `json:"metadata" yaml:"metadata"`
	Data      SecretV1BetaData       `json:"data" yaml:"data"`
	OrgConfig *SecretV1BetaOrgConfig `json:"org_config,omitempty" yaml:"org_config,omitempty"`
}

type SecretV1BetaEnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value,omitempty"`
}

type SecretV1BetaFile struct {
	Path    string `json:"path" yaml:"path"`
	Content string `json:"content" yaml:"content,omitempty"`
}

type SecretV1BetaData struct {
	EnvVars []SecretV1BetaEnvVar `json:"env_vars" yaml:"env_vars"`
	Files   []SecretV1BetaFile   `json:"files" yaml:"files"`
}

type SecretV1BetaMetadata struct {
	Name            string      `json:"name,omitempty" yaml:"name,omitempty"`
	Id              string      `json:"id,omitempty" yaml:"id,omitempty"`
	CreateTime      json.Number `json:"create_time,omitempty,string" yaml:"create_time,omitempty"`
	UpdateTime      json.Number `json:"update_time,omitempty,string" yaml:"update_time,omitempty"`
	ContentIncluded bool        `json:"content_included,omitempty" yaml:"content_included,omitempty"`
}

type SecretV1BetaOrgConfig struct {
	Projects_access string   `json:"projects_access,omitempty" yaml:"projects_access,omitempty"`
	Project_ids     []string `json:"project_ids,omitempty" yaml:"project_ids,omitempty"`
	Debug_access    string   `json:"debug_access,omitempty" yaml:"debug_access,omitempty"`
	Attach_access   string   `json:"attach_access,omitempty" yaml:"attach_access,omitempty"`
}

func NewSecretV1Beta(name string, envVars []SecretV1BetaEnvVar, files []SecretV1BetaFile) SecretV1Beta {
	s := SecretV1Beta{}

	s.setApiVersionAndKind()
	s.Metadata.Name = name
	s.Data.EnvVars = envVars
	s.Data.Files = files

	return s
}

func NewSecretV1BetaFromJson(data []byte) (*SecretV1Beta, error) {
	s := SecretV1Beta{}

	err := json.Unmarshal(data, &s)

	if err != nil {
		return nil, err
	}

	s.setApiVersionAndKind()

	return &s, nil
}

func NewSecretV1BetaFromYaml(data []byte) (*SecretV1Beta, error) {
	s := SecretV1Beta{}

	err := yaml.UnmarshalStrict(data, &s)

	if err != nil {
		return nil, err
	}

	s.setApiVersionAndKind()

	return &s, nil
}

func (s *SecretV1Beta) setApiVersionAndKind() {
	s.ApiVersion = "v1beta"
	s.Kind = "Secret"
}

func (s *SecretV1Beta) ObjectName() string {
	return fmt.Sprintf("Secrets/%s", s.Metadata.Name)
}

func (s *SecretV1Beta) Editable() bool {
	return s.Metadata.ContentIncluded
}

func (s *SecretV1Beta) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SecretV1Beta) ToYaml() ([]byte, error) {
	if !s.Metadata.ContentIncluded {
		notice := []byte(`
# DANGER! Secrets cannot be updated, only replaced. Once the change is applied, the old values will be lost forever.
# Note: You can exit without saving to skip.

`)

		s, err := yaml.Marshal(s)
		if err != nil {
			return nil, err
		}
		return append(notice, s...), nil

	}
	return yaml.Marshal(s)
}
