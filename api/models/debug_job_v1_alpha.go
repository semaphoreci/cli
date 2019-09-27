package models

import "encoding/json"

type DebugJobV1Alpha struct {
	JobId string `json:"job_id,omitempty" yaml:"job_id"`
	Duration int `json:"duration,omitempty,string" yaml:"duration,omitempty"`
}

func NewDebugJobV1Alpha(job_id string, duration int) (*DebugJobV1Alpha) {
	j := DebugJobV1Alpha{}

	j.JobId = job_id
	j.Duration = duration

	return &j
}

func (j *DebugJobV1Alpha) ToJson() ([]byte, error) {
	return json.Marshal(j)
}
