package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v2"
)

type DashboardV1Alpha struct {
	ApiVersion string `json:"apiVersion,omitempty" yaml:"apiVersion"`
	Kind       string `json:"kind,omitempty" yaml:"kind"`
	Metadata   struct {
		Name       string      `json:"name,omitempty"`
		Title      string      `json:"title,omitempty"`
		Id         string      `json:"id,omitempty"`
		CreateTime json.Number `json:"create_time,omitempty,string" yaml:"create_time,omitempty"`
		UpdateTime json.Number `json:"update_time,omitempty,string" yaml:"update_time,omitempty"`
	} `json:"metadata,omitempty"`

	Spec struct {
		Widgets []struct {
			Name    string            `json:"name,omitempty"`
			Type    string            `json:"type,omitempty"`
			Filters map[string]string `json:"filters,omitempty"`
		} `json:"widgets,omitempty"`
	} `json:"spec,omitempty"`
}

func NewDashboardV1Alpha(name string) DashboardV1Alpha {
	d := DashboardV1Alpha{}

	d.ApiVersion = "v1alpha"
	d.Kind = "Dashboard"
	d.Metadata.Name = name

	return d
}

func NewDashboardV1AlphaFromJson(data []byte) (*DashboardV1Alpha, error) {
	d := DashboardV1Alpha{}

	err := json.Unmarshal(data, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func NewDashboardV1AlphaFromYaml(data []byte) (*DashboardV1Alpha, error) {
	d := DashboardV1Alpha{}

	err := yaml.UnmarshalStrict(data, &d)

	if err != nil {
		return nil, err
	}

	return &d, nil
}

func (s *DashboardV1Alpha) ObjectName() string {
	return fmt.Sprintf("Dashboard/%s", s.Metadata.Name)
}

func (s *DashboardV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(s)
}

func (s *DashboardV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(s)
}
