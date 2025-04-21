package client

import (
	"fmt"

	models "github.com/semaphoreci/cli/api/models"
)

type CanvasV2API struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewCanvasV2API() CanvasV2API {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v2")

	return CanvasV2API{
		BaseClient:           baseClient,
		ResourceNamePlural:   "canvases",
		ResourceNameSingular: "canvas",
	}
}

func (c *CanvasV2API) GetCanvas(name string) (*models.CanvasV2, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, name)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewCanvasV2FromJson(body)
}

func (c *CanvasV2API) CreateCanvas(d *models.CanvasV2) (*models.CanvasV2, error) {
	json_body, err := d.ToJson()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize object '%s'", err)
	}

	body, status, err := c.BaseClient.Post(c.ResourceNamePlural, json_body)

	if err != nil {
		return nil, fmt.Errorf("creating %s on Semaphore failed '%s'", c.ResourceNameSingular, err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewCanvasV2FromJson(body)
}
