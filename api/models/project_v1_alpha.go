package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type Scheduler struct {
	Name         string `json:"name"`
	Id           string `json:"id,omitempty"`
	Branch       string `json:"branch"`
	At           string `json:"at"`
	PipelineFile string `json:"pipeline_file" yaml:"pipeline_file"`
	Status       string `json:"status,omitempty" yaml:"status,omitempty"`
}

type Task struct {
	Name         string          `json:"name"`
	Scheduled    bool            `json:"scheduled"`
	Id           string          `json:"id,omitempty"`
	Branch       string          `json:"branch,omitempty"`
	At           string          `json:"at,omitempty"`
	PipelineFile string          `json:"pipeline_file" yaml:"pipeline_file,omitempty"`
	Status       string          `json:"status,omitempty" yaml:"status,omitempty"`
	Parameters   []TaskParameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

type TaskParameter struct {
	Name         string   `json:"name"`
	Required     bool     `json:"required"`
	Description  string   `json:"description,omitempty" yaml:"description,omitempty"`
	DefaultValue string   `json:"default_value,omitempty" yaml:"default_value,omitempty"`
	Options      []string `json:"options,omitempty" yaml:"options,omitempty"`
}

type ForkedPullRequests struct {
	AllowedSecrets      []string `json:"allowed_secrets,omitempty" yaml:"allowed_secrets,omitempty"`
	AllowedContributors []string `json:"allowed_contributors,omitempty" yaml:"allowed_contributors,omitempty"`
}

type Status struct {
	PipelineFiles []PipelineFile `json:"pipeline_files" yaml:"pipeline_files"`
}

type PipelineFile struct {
	Path  string `json:"path"`
	Level string `json:"level"`
}

type Whitelist struct {
	Branches []string `json:"branches,omitempty"`
	Tags     []string `json:"tags,omitempty"`
}

type ProjectV1Alpha struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string `json:"kind,omitempty" yaml:"kind"`
	Metadata   struct {
		Name        string `json:"name,omitempty"`
		Id          string `json:"id,omitempty"`
		Description string `json:"description,omitempty"`
	} `json:"metadata,omitempty"`

	Spec struct {
		Visibility string `json:"visibility,omitempty" yaml:"visibility,omitempty"`
		Repository struct {
			Url                string             `json:"url,omitempty"`
			RunOn              []string           `json:"run_on,omitempty" yaml:"run_on"`
			ForkedPullRequests ForkedPullRequests `json:"forked_pull_requests,omitempty" yaml:"forked_pull_requests,omitempty"`
			PipelineFile       string             `json:"pipeline_file" yaml:"pipeline_file"`
			Status             *Status            `json:"status,omitempty" yaml:"status"`
			Whitelist          Whitelist          `json:"whitelist" yaml:"whitelist"`
			IntegrationType    string             `json:"integration_type" yaml:"integration_type"`
		} `json:"repository,omitempty"`
		Schedulers        []Scheduler `json:"schedulers,omitempty" yaml:"schedulers,omitempty"`
		Tasks             []Task      `json:"tasks,omitempty" yaml:"tasks,omitempty"`
		CustomPermissions *bool       `json:"custom_permissions,omitempty" yaml:"custom_permissions,omitempty"`
		DebugPermissions  []string    `json:"debug_permissions,omitempty" yaml:"debug_permissions,omitempty"`
		AttachPermissions []string    `json:"attach_permissions,omitempty" yaml:"attach_permissions,omitempty"`
	} `json:"spec,omitempty"`
}

func NewProjectV1Alpha(name string) ProjectV1Alpha {
	p := ProjectV1Alpha{}

	p.Metadata.Name = name
	p.setApiVersionAndKind()

	return p
}

func NewProjectV1AlphaFromJson(data []byte) (*ProjectV1Alpha, error) {
	p := ProjectV1Alpha{}

	err := json.Unmarshal(data, &p)

	if err != nil {
		return nil, err
	}

	p.setApiVersionAndKind()

	return &p, nil
}

func NewProjectV1AlphaFromYaml(data []byte) (*ProjectV1Alpha, error) {
	p := ProjectV1Alpha{}

	err := yaml.UnmarshalStrict(data, &p)

	if err != nil {
		return nil, err
	}

	p.setApiVersionAndKind()

	return &p, nil
}

func (p *ProjectV1Alpha) setApiVersionAndKind() {
	p.ApiVersion = "v1alpha"
	p.Kind = "Project"
}

func (p *ProjectV1Alpha) ObjectName() string {
	return fmt.Sprintf("Projects/%s", p.Metadata.Name)
}

func (p *ProjectV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(p)
}

func (p *ProjectV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(p)
}
