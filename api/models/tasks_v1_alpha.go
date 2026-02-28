package models

import (
	"encoding/json"
)

const (
	RunTaskRefBranch = "BRANCH"
	RunTaskRefTag    = "TAG"
)

type TaskV1Alpha struct {
	ID           string            `json:"id" yaml:"id"`
	Name         string            `json:"name" yaml:"name"`
	Description  string            `json:"description,omitempty" yaml:"description,omitempty"`
	ProjectID    string            `json:"project_id" yaml:"project_id"`
	Branch       string            `json:"branch,omitempty" yaml:"branch,omitempty"`
	At           string            `json:"at,omitempty" yaml:"at,omitempty"`
	PipelineFile string            `json:"pipeline_file" yaml:"pipeline_file"`
	RequesterID  string            `json:"requester_id,omitempty" yaml:"requester_id,omitempty"`
	UpdatedAt    string            `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Paused       bool              `json:"paused" yaml:"paused"`
	Suspended    bool              `json:"suspended" yaml:"suspended"`
	Recurring    bool              `json:"recurring" yaml:"recurring"`
	Parameters   map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

type TriggerV1Alpha struct {
	TriggeredAt         string `json:"triggered_at" yaml:"triggered_at"`
	SchedulingStatus    string `json:"scheduling_status" yaml:"scheduling_status"`
	ScheduledWorkflowID string `json:"scheduled_workflow_id,omitempty" yaml:"scheduled_workflow_id,omitempty"`
	Branch              string `json:"branch,omitempty" yaml:"branch,omitempty"`
	PipelineFile        string `json:"pipeline_file,omitempty" yaml:"pipeline_file,omitempty"`
	ErrorDescription    string `json:"error_description,omitempty" yaml:"error_description,omitempty"`
}

type TaskListV1Alpha []TaskV1Alpha

type TaskDescribeV1Alpha struct {
	Schedule TaskV1Alpha      `json:"schedule" yaml:"schedule"`
	Triggers []TriggerV1Alpha `json:"triggers,omitempty" yaml:"triggers,omitempty"`
}

type RunTaskReference struct {
	Type string `json:"type" yaml:"type"`
	Name string `json:"name" yaml:"name"`
}

type RunTaskRequest struct {
	Reference    *RunTaskReference `json:"reference,omitempty" yaml:"reference,omitempty"`
	PipelineFile string            `json:"pipeline_file,omitempty" yaml:"pipeline_file,omitempty"`
	Parameters   map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

type RunTaskResponse struct {
	WorkflowID string `json:"workflow_id" yaml:"workflow_id"`
}

func NewTaskListV1AlphaFromJSON(data []byte) (TaskListV1Alpha, error) {
	var list TaskListV1Alpha
	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func NewTaskDescribeV1AlphaFromJSON(data []byte) (*TaskDescribeV1Alpha, error) {
	t := TaskDescribeV1Alpha{}
	err := json.Unmarshal(data, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func NewRunTaskResponseFromJSON(data []byte) (*RunTaskResponse, error) {
	r := RunTaskResponse{}
	err := json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
