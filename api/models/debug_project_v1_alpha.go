package models

import "encoding/json"

type DebugProjectV1Alpha struct {
	ProjectIdOrName string `json:"project_id_or_name,omitempty" yaml:"project_id_or_name"`
	Duration        int    `json:"duration,omitempty,string" yaml:"duration,omitempty"`
	MachineType     string `json:"machine_type,omitempty" yaml:"machine_type,omitempty"`
}

func NewDebugProjectV1Alpha(project string, duration int, machine string) *DebugProjectV1Alpha {
	j := DebugProjectV1Alpha{}

	j.ProjectIdOrName = project
	j.Duration = duration
	j.MachineType = machine

	return &j
}

func (j *DebugProjectV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(j)
}
