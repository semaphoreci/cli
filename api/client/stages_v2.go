package client

import (
	"fmt"

	models "github.com/semaphoreci/cli/api/models"
)

type StageV2API struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewStageV2API() StageV2API {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v2")

	return StageV2API{
		BaseClient:           baseClient,
		ResourceNamePlural:   "stages",
		ResourceNameSingular: "stage",
	}
}

func (c *StageV2API) List(canvasID string) ([]models.StageV2, error) {
	base := fmt.Sprintf("canvases/%s/%s", canvasID, c.ResourceNamePlural)
	body, status, err := c.BaseClient.List(base)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewStageV2ListFromJson(body)
}

func (c *StageV2API) GetStage(canvasID, name string) (*models.StageV2, error) {
	base := fmt.Sprintf("canvases/%s/%s", canvasID, c.ResourceNamePlural)
	body, status, err := c.BaseClient.Get(base, name)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewStageV2FromJson(body)
}

func (c *StageV2API) CreateStage(canvasID string, d *models.StageV2) (*models.StageV2, error) {
	base := fmt.Sprintf("canvases/%s/%s", canvasID, c.ResourceNamePlural)
	json_body, err := d.ToJson()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize object '%s'", err)
	}

	body, status, err := c.BaseClient.Post(base, json_body)

	if err != nil {
		return nil, fmt.Errorf("creating %s on Semaphore failed '%s'", c.ResourceNameSingular, err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewStageV2FromJson(body)
}
