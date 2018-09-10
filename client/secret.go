package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v2"
)

type Secret struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind"`

	Metadata struct {
		Name       string `json:"name,omitempty" yaml:"name,omitempty"`
		Id         string `json:"id,omitempty" yaml:"id,omitempty"`
		CreateTime string `json:"create_time,omitempty" yaml:"create_time,omitempty"`
		UpdateTime string `json:"update_time,omitempty" yaml:"update_time,omitempty"`
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

func InitSecret(name string) Secret {
	s := Secret{}

	s.ApiVersion = "v1beta"
	s.Kind = "Secret"
	s.Metadata.Name = name

	return s
}

func InitSecretFromYaml(data []byte) (Secret, error) {
	s := Secret{}

	err := yaml.UnmarshalStrict(data, &s)

	if err != nil {
		return s, err
	}

	if s.ApiVersion == "" {
		s.ApiVersion = "v1beta"
	}

	if s.Kind == "" {
		s.Kind = "Secret"
	}

	return s, nil
}

func InitSecretFromJson(data []byte) (Secret, error) {
	s := Secret{}

	err := json.Unmarshal(data, &s)

	if s.ApiVersion == "" {
		s.ApiVersion = "v1beta"
	}

	if s.Kind == "" {
		s.Kind = "Secret"
	}

	return s, err
}

func GetSecret(name string) (*Secret, error) {
	c := FromConfig()
	c.SetApiVersion("v1beta")

	body, status, err := c.Get("secrets", name)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	s, err := InitSecretFromJson(body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to deserialize secret object '%s'", err))
	}

	return &s, nil
}

func DeleteSecret(name string) error {
	c := FromConfig()
	c.SetApiVersion("v1beta")

	body, status, err := c.Delete("secrets", name)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return nil
}

func (s *Secret) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *Secret) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}

func (s *Secret) ObjectName() string {
	return fmt.Sprintf("Secrets/%s", s.Metadata.Name)
}

func (s *Secret) Create() error {
	c := FromConfig()
	c.SetApiVersion("v1beta")

	json_body, err := s.ToJson()

	if err != nil {
		return errors.New(fmt.Sprintf("failed to serialize secret object '%s'", err))
	}

	body, status, err := c.Post("secrets", json_body)

	if err != nil {
		return errors.New(fmt.Sprintf("creating secret on Semaphore failed '%s'", err))
	}

	if status != 200 {
		return errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return nil
}

func (s *Secret) Update() error {
	c := FromConfig()
	c.SetApiVersion("v1beta")

	json_body, err := s.ToJson()

	if err != nil {
		return errors.New(fmt.Sprintf("failed to serialize secret object '%s'", err))
	}

	body, status, err := c.Patch("secrets", s.Metadata.Name, json_body)

	if err != nil {
		return errors.New(fmt.Sprintf("updating secret on Semaphore failed '%s'", err))
	}

	if status != 200 {
		return errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return nil
}
