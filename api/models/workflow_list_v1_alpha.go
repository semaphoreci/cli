package models

import "encoding/json"

type WorkflowListV1Alpha struct {
	Workflow []WorkflowV1Alpha `json:"workflows" yaml:"projects"`
}

func NewWorkflowListV1AlphaFromJson(data []byte) (*WorkflowListV1Alpha, error) {
	list := []WorkflowV1Alpha{}

	err := json.Unmarshal(data, &list)

	if err != nil {
		return nil, err
	}

	return &WorkflowListV1Alpha{Workflow: list}, nil
}
