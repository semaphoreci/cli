package client

import (
	"fmt"
	"net/url"

	models "github.com/semaphoreci/cli/api/models"
)

type AgentApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewAgentApiV1AlphaApi() AgentApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return AgentApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "self_hosted_agents",
		ResourceNameSingular: "self_hosted_agent",
	}
}

func (c *AgentApiV1AlphaApi) ListAgents(agentType string, next int64) (*models.AgentListV1Alpha, error) {
	query := url.Values{}
	query.Add("count", "200")

	if agentType != "" {
		query.Add("agent_type", agentType)
	}

	if next > 0 {
		query.Add("next", fmt.Sprintf("%d", next))
	}

	body, status, err := c.BaseClient.ListWithParams(c.ResourceNamePlural, query)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewAgentListV1AlphaFromJson(body)
}

func (c *AgentApiV1AlphaApi) GetAgent(name string) (*models.AgentV1Alpha, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, name)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewAgentV1AlphaFromJson(body)
}
