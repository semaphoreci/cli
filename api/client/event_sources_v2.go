package client

import (
	"fmt"

	models "github.com/semaphoreci/cli/api/models"
)

type EventSourceV2API struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewEventSourceV2API() EventSourceV2API {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v2")

	return EventSourceV2API{
		BaseClient:           baseClient,
		ResourceNamePlural:   "sources",
		ResourceNameSingular: "source",
	}
}

func (c *EventSourceV2API) List(canvasID string) ([]models.EventSourceV2, error) {
	base := fmt.Sprintf("canvases/%s/%s", canvasID, c.ResourceNamePlural)
	body, status, err := c.BaseClient.List(base)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewEventSourceV2ListFromJson(body)
}

func (c *EventSourceV2API) CreateEventSource(canvasID string, d *models.EventSourceV2) (*models.EventSourceV2, error) {
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

	return models.NewEventSourceV2FromJson(body)
}

func (c *EventSourceV2API) GetEventSource(canvasID string, name string) (*models.EventSourceV2, error) {
	base := fmt.Sprintf("canvases/%s/%s", canvasID, c.ResourceNamePlural)
	body, status, err := c.BaseClient.Get(base, name)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewEventSourceV2FromJson(body)
}
