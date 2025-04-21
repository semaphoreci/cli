package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

const KindEventSource = "EventSource"

type EventSourceV2 struct {
	ApiVersion string                `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                `json:"kind" yaml:"kind"`
	Metadata   EventSourceV2Metadata `json:"metadata" yaml:"metadata"`
	Spec       EventSourceV2Spec     `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type EventSourceV2Metadata struct {
	ID                  string                  `json:"id,omitempty" yaml:"id"`
	Name                string                  `json:"name" yaml:"name"`
	Organization        *OrganizationMetadataV2 `json:"organization,omitempty" yaml:"organization"`
	Canvas              *CanvasMetadataV2       `json:"canvas" yaml:"canvas"`
	Timeline            *TimelineMetadataV2     `json:"timeline,omitempty" yaml:"timeline"`
	EventSourceV2Status *EventSourceV2Status    `json:"status,omitempty" yaml:"status"`
}

type EventSourceV2Status struct {
	Key string `json:"key,omitempty" yaml:"key"`
}

type EventSourceV2Spec struct{}

func NewEventSourceV2(name string) EventSourceV2 {
	e := EventSourceV2{}
	e.Metadata.Name = name
	e.setApiVersionAndKind()
	return e
}

func NewEventSourceV2FromJson(data []byte) (*EventSourceV2, error) {
	e := EventSourceV2{}

	err := json.Unmarshal(data, &e)
	if err != nil {
		return nil, err
	}

	e.setApiVersionAndKind()
	return &e, nil
}

func NewEventSourceV2ListFromJson(data []byte) ([]EventSourceV2, error) {
	e := []EventSourceV2{}
	err := json.Unmarshal(data, &e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func NewEventSourceV2FromYaml(data []byte) (*EventSourceV2, error) {
	e := EventSourceV2{}

	err := yaml.UnmarshalStrict(data, &e)
	if err != nil {
		return nil, err
	}

	e.setApiVersionAndKind()

	return &e, nil
}

func (s *EventSourceV2) setApiVersionAndKind() {
	s.ApiVersion = "v2"
	s.Kind = KindEventSource
}

func (s *EventSourceV2) ObjectName() string {
	return fmt.Sprintf("%s/%s", KindEventSource, s.Metadata.Name)
}

func (s *EventSourceV2) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *EventSourceV2) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}
