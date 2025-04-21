package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

const KindCanvas = "Canvas"

type OrganizationMetadataV2 struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

type CanvasMetadataV2 struct {
	ID   string `json:"id" yaml:"id"`
	Name string `json:"name" yaml:"name"`
}

type TimelineMetadataV2 struct {
	CreatedAt string `json:"created_at" yaml:"created_at"`
}

type CanvasV2 struct {
	ApiVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   CanvasV2Metadata `json:"metadata" yaml:"metadata"`
	Spec       CanvasV2Spec     `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type CanvasV2Metadata struct {
	ID           string                 `json:"id,omitempty" yaml:"id,omitempty"`
	Name         string                 `json:"name" yaml:"name"`
	Organization OrganizationMetadataV2 `json:"organization,omitempty" yaml:"organization,omitempty"`
	Timeline     *TimelineMetadataV2    `json:"timeline,omitempty" yaml:"timeline,omitempty"`
}

type CanvasV2Spec struct{}

func NewCanvasV2(name string) CanvasV2 {
	a := CanvasV2{}
	a.Metadata.Name = name
	a.setApiVersionAndKind()
	return a
}

func NewCanvasV2FromJson(data []byte) (*CanvasV2, error) {
	a := CanvasV2{}

	err := json.Unmarshal(data, &a)
	if err != nil {
		return nil, err
	}

	a.setApiVersionAndKind()
	return &a, nil
}

func NewCanvasV2FromYaml(data []byte) (*CanvasV2, error) {
	a := CanvasV2{}

	err := yaml.UnmarshalStrict(data, &a)
	if err != nil {
		return nil, err
	}

	a.setApiVersionAndKind()

	return &a, nil
}

func (s *CanvasV2) setApiVersionAndKind() {
	s.ApiVersion = "v2"
	s.Kind = KindCanvas
}

func (s *CanvasV2) ObjectName() string {
	return fmt.Sprintf("%s/%s", KindCanvas, s.Metadata.Name)
}

func (s *CanvasV2) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *CanvasV2) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}
