package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type SecretV1Beta struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string `json:"kind,omitempty" yaml:"kind"`

	Metadata SecretV1BetaMetadata `json:"metadata" yaml:"metadata"`
	Data     SecretV1BetaData     `json:"data" yaml:"data"`
}

type SecretV1BetaEnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

type SecretV1BetaFile struct {
	Path    string `json:"path" yaml:"path"`
	Content string `json:"content" yaml:"content"`
}

type SecretV1BetaData struct {
	EnvVars []SecretV1BetaEnvVar `json:"env_vars" yaml:"env_vars"`
	Files   []SecretV1BetaFile   `json:"files" yaml:"files"`
}

type SecretV1BetaMetadata struct {
	Name       string      `json:"name,omitempty" yaml:"name,omitempty"`
	Id         string      `json:"id,omitempty" yaml:"id,omitempty"`
	CreateTime json.Number `json:"create_time,omitempty,string" yaml:"create_time,omitempty"`
	UpdateTime json.Number `json:"update_time,omitempty,string" yaml:"update_time,omitempty"`
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

func (s *SecretV1Beta) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *SecretV1Beta) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}
