package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ghodss/yaml"
)

type Secret struct {
	Metadata struct {
		Name       string `json:"name,omitempty"`
		Id         string `json:"id,omitempty"`
		CreateTime string `json:"create_time,omitempty"`
		UpdateTime string `json:"update_time,omitempty"`
	} `json:"metadata"`

	Data struct {
		EnvVars []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"env_vars"`

		Files []struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		} `json:"files"`
	} `json:"data"`
}

func InitSecret(name string) Secret {
	s := Secret{}

	s.Metadata.Name = name

	return s
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

	s := Secret{}
	err = json.Unmarshal(body, &s)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to deserialize secret object '%s'", err))
	}

	return &s, nil
}

func InitSecretFromYaml(data []byte) (Secret, error) {
	s := Secret{}

	err := yaml.Unmarshal(data, &s)

	if err != nil {
		return s, err
	}

	return s, nil
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
