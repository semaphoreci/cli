package models

import "encoding/json"

type SecretListV1Beta struct {
	Secrets []SecretV1Beta `json:"secrets" yaml:"secrets"`
}

func NewSecretListV1BetaFromJson(data []byte) (*SecretListV1Beta, error) {
	list := SecretListV1Beta{}

	err := json.Unmarshal(data, &list)

	if err != nil {
		return nil, err
	}

	for _, s := range list.Secrets {
		if s.ApiVersion == "" {
			s.ApiVersion = "v1beta"
		}

		if s.Kind == "" {
			s.Kind = "Secret"
		}
	}

	return &list, nil
}
