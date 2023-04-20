package models

import (
	"encoding/json"

	yaml "gopkg.in/yaml.v2"
)

type TroubleshootV1Alpha struct {
	Workflow map[string]interface{} `json:"workflow,omitempty" yaml:"workflow,omitempty"`
	Project  map[string]interface{} `json:"project,omitempty" yaml:"project,omitempty"`
	Pipeline map[string]interface{} `json:"pipeline,omitempty" yaml:"pipeline,omitempty"`
	Job      map[string]interface{} `json:"job,omitempty" yaml:"job,omitempty"`
	Block    map[string]interface{} `json:"block,omitempty" yaml:"block,omitempty"`
}

func NewTroubleshootV1AlphaFromJson(data []byte) (*TroubleshootV1Alpha, error) {
	t := TroubleshootV1Alpha{}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (j *TroubleshootV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(j)
}

func (j *TroubleshootV1Alpha) ToYaml() ([]byte, error) {
	return yaml.Marshal(j)
}
