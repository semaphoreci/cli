package client

import (
	"errors"
	"fmt"

	models "github.com/semaphoreci/cli/api/models"
)

type AgentTypeApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewAgentTypeApiV1AlphaApi() AgentTypeApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return AgentTypeApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "self_hosted_agent_types",
		ResourceNameSingular: "self_hosted_agent_type",
	}
}

func (c *AgentTypeApiV1AlphaApi) ListAgentTypes() (*models.AgentTypeListV1Alpha, error) {
	body, status, err := c.BaseClient.List(c.ResourceNamePlural)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewAgentTypeListV1AlphaFromJson(body)
}

func (c *AgentTypeApiV1AlphaApi) GetAgentType(name string) (*models.AgentTypeV1Alpha, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, name)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewAgentTypeV1AlphaFromJson(body)
}

func (c *AgentTypeApiV1AlphaApi) DeleteAgentType(name string) error {
	body, status, err := c.BaseClient.Delete(c.ResourceNamePlural, name)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return nil
}

func (c *AgentTypeApiV1AlphaApi) CreateAgentType(d *models.AgentTypeV1Alpha) (*models.AgentTypeV1Alpha, error) {
	json_body, err := d.ToJson()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to serialize object '%s'", err))
	}

	body, status, err := c.BaseClient.Post(c.ResourceNamePlural, json_body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("creating %s on Semaphore failed '%s'", c.ResourceNameSingular, err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewAgentTypeV1AlphaFromJson(body)
}

func (c *AgentTypeApiV1AlphaApi) UpdateAgentType(d *models.AgentTypeV1Alpha) (*models.AgentTypeV1Alpha, error) {
	json_body, err := d.ToJson()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to serialize object '%s'", err))
	}

	body, status, err := c.BaseClient.Patch(c.ResourceNamePlural, d.Metadata.Name, json_body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("updating %s on Semaphore failed '%s'", c.ResourceNameSingular, err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewAgentTypeV1AlphaFromJson(body)
}
