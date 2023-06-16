package models

import "encoding/json"

type AgentListV1Alpha struct {
	Agents []AgentV1Alpha `json:"agents" yaml:"agents"`
	Cursor string         `json:"cursor" yaml:"cursor"`
}

func NewAgentListV1AlphaFromJson(data []byte) (*AgentListV1Alpha, error) {
	list := AgentListV1Alpha{}

	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	for _, s := range list.Agents {
		if s.ApiVersion == "" {
			s.ApiVersion = "v1alpha"
		}

		if s.Kind == "" {
			s.Kind = "SelfHostedAgent"
		}
	}

	return &list, nil
}
