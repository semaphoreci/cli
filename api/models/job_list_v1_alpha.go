package models

import (
	"encoding/json"
)

type JobListV1Alpha struct {
	Jobs []JobV1Alpha `json:"jobs" yaml:"jobs"`
}

func NewJobListV1AlphaFromJson(data []byte) (*JobListV1Alpha, error) {
	list := JobListV1Alpha{}

	err := json.Unmarshal(data, &list)

	if err != nil {
		return nil, err
	}

	for _, j := range list.Jobs {
		if j.ApiVersion == "" {
			j.ApiVersion = "v1alpha"
		}

		if j.Kind == "" {
			j.Kind = "Job"
		}
	}

	return &list, nil
}
