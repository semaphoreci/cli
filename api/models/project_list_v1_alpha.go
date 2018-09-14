package models

import "encoding/json"

type ProjectListV1Alpha struct {
	Projects []ProjectListV1Alpha `json:"projects" yaml:"projects"`
}

func NewProjectListV1AlphaFromJson(data []byte) (*ProjectListV1Alpha, error) {
	list := []ProjectV1Alpha{}

	err := json.Unmarshal(data, &list)

	if err != nil {
		return nil, err
	}

	for _, s := range list.Dashboards {
		if s.ApiVersion == "" {
			s.ApiVersion = "v1alpha"
		}

		if s.Kind == "" {
			s.Kind = "Project"
		}
	}

	return &ProjectListV1Alpha{Projects: list}, nil
}
