package models

import "encoding/json"

type ProjectListV1Alpha struct {
	Projects []ProjectV1Alpha `json:"projects" yaml:"projects"`
}

func NewProjectListV1AlphaFromJson(data []byte) (*ProjectListV1Alpha, error) {
	list := []ProjectV1Alpha{}

	err := json.Unmarshal(data, &list)

	if err != nil {
		return nil, err
	}

	for _, p := range list {
		if p.ApiVersion == "" {
			p.ApiVersion = "v1alpha"
		}

		if p.Kind == "" {
			p.Kind = "Project"
		}
	}

	return &ProjectListV1Alpha{Projects: list}, nil
}
