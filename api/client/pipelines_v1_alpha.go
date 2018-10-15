package client

import (
	"errors"
	"fmt"

	models "github.com/semaphoreci/cli/api/models"
)

type PipelinesApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewPipelinesV1AlphaApi() PipelinesApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return PipelinesApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "pipelines",
		ResourceNameSingular: "pipeline",
	}
}

func (c *PipelinesApiV1AlphaApi) DescribePpl(id string) (*models.PipelineV1Alpha, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, id)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewPipelineV1AlphaFromJson(body)
}
