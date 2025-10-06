package models

import (
	"encoding/json"
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Reference struct {
	Type string `json:"type" yaml:"type"`
	Name string `json:"name" yaml:"name"`
}

type Scheduler struct {
	Name         string     `json:"name"`
	Id           string     `json:"id,omitempty"`
	Branch       string     `json:"branch,omitempty" yaml:"branch,omitempty"` // deprecated: use Reference instead
	Reference    *Reference `json:"reference,omitempty" yaml:"reference,omitempty"`
	At           string     `json:"at"`
	PipelineFile string     `json:"pipeline_file" yaml:"pipeline_file"`
	Status       string     `json:"status,omitempty" yaml:"status,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for backward compatibility
func (s *Scheduler) UnmarshalJSON(data []byte) error {
	type Alias Scheduler
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// If we have a reference field from the API, use it and clear branch
	// to avoid duplication in the output
	if s.Reference != nil {
		s.Branch = ""
		return nil
	}

	// For schedulers, always create a reference from branch if not provided
	// This ensures consistent output format
	if s.Branch != "" {
		s.Reference = referenceFromBranch(s.Branch)
		// Clear the branch field since we now have a reference
		s.Branch = ""
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for backward compatibility
func (s *Scheduler) MarshalJSON() ([]byte, error) {
	type Alias Scheduler

	// Create a copy to avoid mutating the original struct
	temp := *s

	// If we have a Reference but no Branch, convert Reference to Branch
	// for backward compatibility and omit the Reference field
	if s.Reference != nil && s.Branch == "" {
		if s.Reference.Type == "branch" {
			temp.Branch = s.Reference.Name
		} else if s.Reference.Type == "tag" {
			temp.Branch = refTagsPrefix + s.Reference.Name
		}
		// Don't include Reference in output to avoid confusion on re-parse
		temp.Reference = nil
	}

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&temp),
	})
}

type Task struct {
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Scheduled    bool            `json:"scheduled"`
	Id           string          `json:"id,omitempty"`
	Branch       string          `json:"branch,omitempty"` // deprecated: use Reference instead
	Reference    *Reference      `json:"reference,omitempty" yaml:"reference,omitempty"`
	At           string          `json:"at,omitempty"`
	PipelineFile string          `json:"pipeline_file" yaml:"pipeline_file,omitempty"`
	Status       string          `json:"status,omitempty" yaml:"status,omitempty"`
	Parameters   []TaskParameter `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for backward compatibility
func (t *Task) UnmarshalJSON(data []byte) error {
	type Alias Task
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// If we have a reference field from the API, use it and keep branch
	// for backward compatibility (branch field is preserved when both exist)
	if t.Reference != nil {
		return nil
	}

	// For tasks, always create a reference from branch if not provided
	// This ensures consistent output format
	if t.Branch != "" {
		t.Reference = referenceFromBranch(t.Branch)
		// Clear the branch field since we now have a reference
		t.Branch = ""
	}

	return nil
}

// MarshalJSON implements custom JSON marshaling for backward compatibility
func (t *Task) MarshalJSON() ([]byte, error) {
	type Alias Task

	// Create a copy to avoid mutating the original struct
	temp := *t

	// If we have a Reference but no Branch, convert Reference to Branch
	// for backward compatibility and omit the Reference field
	if t.Reference != nil && t.Branch == "" {
		if t.Reference.Type == "branch" {
			temp.Branch = t.Reference.Name
		} else if t.Reference.Type == "tag" {
			temp.Branch = refTagsPrefix + t.Reference.Name
		}
		// Don't include Reference in output to avoid confusion on re-parse
		temp.Reference = nil
	}

	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(&temp),
	})
}

// MarshalYAML implements custom YAML marshaling for backward compatibility
func (t *Task) MarshalYAML() (interface{}, error) {
	type Alias Task

	// Create a copy to avoid mutating the original struct
	temp := *t

	// If we have a Reference but no Branch, convert Reference to Branch
	// for backward compatibility and omit the Reference field
	if t.Reference != nil && t.Branch == "" {
		if t.Reference.Type == "branch" {
			temp.Branch = t.Reference.Name
		} else if t.Reference.Type == "tag" {
			temp.Branch = refTagsPrefix + t.Reference.Name
		}
		// Don't include Reference in output to avoid confusion on re-parse
		temp.Reference = nil
	}

	return &struct {
		*Alias
	}{
		Alias: (*Alias)(&temp),
	}, nil
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

const refTagsPrefix = "refs/tags/"

// referenceFromBranch converts a legacy branch string into a Reference object.
// Returns nil if the branch is empty.
func referenceFromBranch(branch string) *Reference {
	if branch == "" {
		return nil
	}

	if strings.HasPrefix(branch, refTagsPrefix) {
		return &Reference{
			Type: "tag",
			Name: branch[len(refTagsPrefix):],
		}
	}

	return &Reference{
		Type: "branch",
		Name: branch,
	}
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
