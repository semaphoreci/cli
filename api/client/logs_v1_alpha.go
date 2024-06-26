package client

import (
	"fmt"

	models "github.com/semaphoreci/cli/api/models"
)

type LogsApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewLogsV1AlphaApi() LogsApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return LogsApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "logs",
		ResourceNameSingular: "logs",
	}
}

func (c *LogsApiV1AlphaApi) Get(jobID string) (*models.LogsV1Alpha, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, jobID)
	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewLogsV1AlphaFromJson(body)
}
