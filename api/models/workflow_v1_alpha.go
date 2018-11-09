package models

type WorkflowV1Alpha struct {
	Id           string `json:"wf_id,omitempty" yaml:"id,omitempty"`
	InitialPplId string `json:"initial_ppl_id,omitempty" yaml:"initial_ppl_id,omitempty"`
	BranchName   string `json:"branch_name,omitempty" yaml:"branch_name,omitempty"`
	CreatedAt struct {
		Seconds int64 `json:"seconds"`
	} `json:"created_at"`
}
