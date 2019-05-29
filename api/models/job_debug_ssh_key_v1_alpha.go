package models

import "encoding/json"

type JobDebugSSHKeyV1Alpha struct {
	Key string `json:"key,omitempty" yaml:"key"`
}

func NewJobDebugSSHKeyV1AlphaFromJSON(data []byte) (*JobDebugSSHKeyV1Alpha, error) {
	key := JobDebugSSHKeyV1Alpha{}

	err := json.Unmarshal(data, &key)

	if err != nil {
		return nil, err
	}

	return &key, nil
}
