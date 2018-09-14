package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type SecretV1Beta struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`

	Metadata struct {
		Name       string `json:"name,omitempty" yaml:"name,omitempty"`
		Id         string `json:"id,omitempty" yaml:"id,omitempty"`
		CreateTime int64  `json:"create_time,omitempty,string" yaml:"create_time,omitempty"`
		UpdateTime int64  `json:"update_time,omitempty,string" yaml:"update_time,omitempty"`
	} `json:"metadata" yaml:"metadata"`

	Data struct {
		EnvVars []struct {
			Name  string `json:"name" yaml:"name"`
			Value string `json:"value" yaml:"value"`
		} `json:"env_vars" yaml:"env_vars"`

		Files []struct {
			Path    string `json:"path" yaml:"path"`
			Content string `json:"content" yaml:"content"`
		} `json:"files" yaml: "files"`
	} `json:"data" yaml: "data"`
}

func NewSecretV1Beta(name string) SecretV1Beta {
	d := SecretV1Beta{}

	d.ApiVersion = "v1beta"
	d.Kind = "Secret"
	d.Metadata.Name = name

	return d
}

func NewSecretV1BetaFromJson(data []byte) (*SecretV1Beta, error) {
	d := SecretV1Beta{}

	err := yaml.UnmarshalStrict(data, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func NewSecretV1BetaFromYaml(data []byte) (*SecretV1Beta, error) {
	d := SecretV1Beta{}

	err := yaml.UnmarshalStrict(data, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
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
