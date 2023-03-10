package models

import "encoding/json"

type ProjectSecretListV1 struct {
	Secrets []ProjectSecretV1 `json:"secrets" yaml:"secrets"`
}

func NewProjectSecretListV1FromJson(data []byte) (*ProjectSecretListV1, error) {
	list := ProjectSecretListV1{}

	err := json.Unmarshal(data, &list)

	if err != nil {
		return nil, err
	}

	for _, s := range list.Secrets {
		if s.ApiVersion == "" {
			s.ApiVersion = "v1"
		}

		if s.Kind == "" {
			s.Kind = "Secret"
		}
	}

	return &list, nil
}
