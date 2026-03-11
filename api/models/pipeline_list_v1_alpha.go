package models

import "encoding/json"

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

func (p *PipelinesListV1Alpha) UnmarshalJSON(data []byte) error {
	var items []PplListElemV1Alpha
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	*p = items
	return nil
}

func (p PipelinesListV1Alpha) MarshalJSON() ([]byte, error) {
	return json.Marshal([]PplListElemV1Alpha(p))
}
