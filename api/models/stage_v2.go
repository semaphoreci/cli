package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

const KindStage = "Stage"

type StageV2 struct {
	ApiVersion string          `json:"apiVersion" yaml:"apiVersion"`
	Kind       string          `json:"kind" yaml:"kind"`
	Metadata   StageV2Metadata `json:"metadata" yaml:"metadata"`
	Spec       StageV2Spec     `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type StageV2Metadata struct {
	ID           string                  `json:"id,omitempty" yaml:"id"`
	Name         string                  `json:"name" yaml:"name"`
	Organization *OrganizationMetadataV2 `json:"organization,omitempty" yaml:"organization"`
	Canvas       *CanvasMetadataV2       `json:"canvas,omitempty" yaml:"canvas"`
	Timeline     *TimelineMetadataV2     `json:"timeline,omitempty" yaml:"timeline"`
}

type StageV2Spec struct {
	ApprovalRequired bool                `json:"approval_required" yaml:"approval_required"`
	Connections      []StageV2Connection `json:"connections" yaml:"connections"`
	Run              *StageV2RunTemplate `json:"run,omitempty" yaml:"run"`
}

type StageV2Connection struct {
	Type           string                    `json:"type" yaml:"type"`
	Name           string                    `json:"name" yaml:"name"`
	FilterOperator string                    `json:"filter_operator,omitempty" yaml:"filter_operator"`
	Filters        []StageV2ConnectionFilter `json:"filters,omitempty" yaml:"filters"`
}

type StageV2ConnectionFilter struct {
	Type string             `json:"type" yaml:"type"`
	Data *StageV2DataFilter `json:"data,omitempty" yaml:"data,omitempty"`
}

type StageV2DataFilter struct {
	Expression string `json:"expression" yaml:"expression"`
}

type StageV2RunTemplate struct {
	Type      string                       `json:"type" yaml:"type"`
	Semaphore *StageV2RunSemaphoreTemplate `json:"semaphore,omitempty" yaml:"semaphore,omitempty"`
}

type StageV2RunSemaphoreTemplate struct {
	ProjectID    string              `json:"project_id" yaml:"project_id"`
	TaskID       string              `json:"task_id,omitempty" yaml:"task_id"`
	Branch       string              `json:"branch" yaml:"branch"`
	PipelineFile string              `json:"pipeline_file" yaml:"pipeline_file"`
	Parameters   []TemplateParameter `json:"parameters,omitempty" yaml:"parameters"`
}

type TemplateParameter struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

func NewStageV2(name string) StageV2 {
	s := StageV2{}
	s.Metadata.Name = name
	s.setApiVersionAndKind()
	return s
}

func NewStageV2FromJson(data []byte) (*StageV2, error) {
	s := StageV2{}

	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}

	s.setApiVersionAndKind()
	return &s, nil
}

func NewStageV2ListFromJson(data []byte) ([]StageV2, error) {
	s := []StageV2{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func NewStageV2FromYaml(data []byte) (*StageV2, error) {
	s := StageV2{}

	err := yaml.UnmarshalStrict(data, &s)
	if err != nil {
		return nil, err
	}

	s.setApiVersionAndKind()

	return &s, nil
}

func (s *StageV2) setApiVersionAndKind() {
	s.ApiVersion = "v2"
	s.Kind = KindStage
}

func (s *StageV2) ObjectName() string {
	return fmt.Sprintf("%s/%s", KindStage, s.Metadata.Name)
}

func (s *StageV2) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *StageV2) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}
