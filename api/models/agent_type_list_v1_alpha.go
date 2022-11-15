package models

import "encoding/json"

type AgentTypeListV1Alpha struct {
	AgentTypes []AgentTypeV1Alpha `json:"agent_types" yaml:"agent_types"`
}

func NewAgentTypeListV1AlphaFromJson(data []byte) (*AgentTypeListV1Alpha, error) {
	list := AgentTypeListV1Alpha{}

	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, err
	}

	for _, s := range list.AgentTypes {
		if s.ApiVersion == "" {
			s.ApiVersion = "v1alpha"
		}

		if s.Kind == "" {
			s.Kind = "SelfHostedAgentType"
		}
	}

	return &list, nil
}
