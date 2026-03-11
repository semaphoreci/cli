package models

type PipelinesListV1Alpha []PplListElemV1Alpha

type PplListElemV1Alpha struct {
	Id        string `json:"ppl_id"`
	Name      string `json:"name"`
	State     string `json:"state"`
	CreatedAt struct {
		Seconds int64 `json:"seconds"`
	} `json:"created_at"`
	Label string `json:"branch_name"`
}
