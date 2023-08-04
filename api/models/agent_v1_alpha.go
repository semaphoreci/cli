package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

type AgentV1Alpha struct {
	ApiVersion string               `json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string               `json:"kind,omitempty" yaml:"kind"`
	Metadata   AgentV1AlphaMetadata `json:"metadata" yaml:"metadata"`
	Status     AgentV1AlphaStatus   `json:"status" yaml:"status"`
}

type AgentV1AlphaMetadata struct {
	Name        string      `json:"name,omitempty" yaml:"name,omitempty"`
	Type        string      `json:"type,omitempty" yaml:"type,omitempty"`
	ConnectedAt json.Number `json:"connected_at,omitempty" yaml:"connected_at,omitempty"`
	DisabledAt  json.Number `json:"disabled_at,omitempty" yaml:"disabled_at,omitempty"`
	Version     string      `json:"version,omitempty" yaml:"version,omitempty"`
	OS          string      `json:"os,omitempty" yaml:"os,omitempty"`
	Arch        string      `json:"arch,omitempty" yaml:"arch,omitempty"`
	Hostname    string      `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	IPAddress   string      `json:"ip_address,omitempty" yaml:"ip_address,omitempty"`
	PID         json.Number `json:"pid,omitempty" yaml:"pid,omitempty"`
}

type AgentV1AlphaStatus struct {
	State string `json:"state,omitempty" yaml:"state,omitempty"`
}

func NewAgentV1Alpha(name string) AgentV1Alpha {
	a := AgentV1Alpha{}
	a.Metadata.Name = name
	a.setApiVersionAndKind()
	return a
}

func NewAgentV1AlphaFromJson(data []byte) (*AgentV1Alpha, error) {
	a := AgentV1Alpha{}

	err := json.Unmarshal(data, &a)
	if err != nil {
		return nil, err
	}

	a.setApiVersionAndKind()
	return &a, nil
}

func NewAgentV1AlphaFromYaml(data []byte) (*AgentV1Alpha, error) {
	a := AgentV1Alpha{}

	err := yaml.UnmarshalStrict(data, &a)
	if err != nil {
		return nil, err
	}

	a.setApiVersionAndKind()
	return &a, nil
}

func (s *AgentV1Alpha) setApiVersionAndKind() {
	s.ApiVersion = "v1alpha"
	s.Kind = "SelfHostedAgent"
}

func (s *AgentV1Alpha) ObjectName() string {
	return fmt.Sprintf("SelfHostedAgent/%s", s.Metadata.Name)
}

func (s *AgentV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *AgentV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}
