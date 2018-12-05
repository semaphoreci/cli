package models

import "encoding/json"

type WorkflowV1Alpha struct {
	Id           string `json:"wf_id,omitempty" yaml:"id,omitempty"`
	InitialPplId string `json:"initial_ppl_id,omitempty" yaml:"initial_ppl_id,omitempty"`
	BranchName   string `json:"branch_name,omitempty" yaml:"branch_name,omitempty"`
	CreatedAt    struct {
		Seconds int64 `json:"seconds"`
	} `json:"created_at"`
}

type WorkflowSnapshotResponseV1Alpha struct {
	WfID  string `json:"wf_id,omitempty" yaml:"id,omitempty"`
	PplID string `json:"ppl_id,omitempty" yaml:"initial_ppl_id,omitempty"`
}

func NewWorkflowSnapshotResponseV1AlphaFromJson(data []byte) (*WorkflowSnapshotResponseV1Alpha, error) {
	j := WorkflowSnapshotResponseV1Alpha{}

	err := json.Unmarshal(data, &j)

	if err != nil {
		return nil, err
	}

	return &j, nil
}
