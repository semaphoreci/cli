package models

type WorkflowV1Alpha struct {
	Id           string `json:"wf_id,omitempty" yaml:"id,omitempty"`
	InitialPplId string `json:"initial_ppl_id,omitempty" yaml:"initial_ppl_id,omitempty"`
	BranchName   string `json:"branch_name,omitempty" yaml:"branch_name,omitempty"`
}

type WorkflowV1AlphaSnapshotRequest struct {
	ProjectID       string `json:"project_id"`
	SnapshorArchive string `json:"snapshot_archive"`
	RequestToken    string `json:"request_token"`
}

func NewWorkflowV1AlphaSnapshotRequest(projectID, snapshorArchive, reqToken string) (*WorkflowV1AlphaSnapshotRequest, error) {
	j := WorkflowV1AlphaSnapshotRequest{}
	j.ProjectID = projectID
	j.SnapshorArchive = snapshorArchive
	j.RequestToken = reqToken
	return &j, nil
}
