package client

import (
	"fmt"

	models "github.com/semaphoreci/cli/api/models"
)

type DashboardApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewDashboardV1AlphaApi() DashboardApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return DashboardApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "dashboards",
		ResourceNameSingular: "dashboard",
	}
}

func (c *DashboardApiV1AlphaApi) ListDashboards() (*models.DashboardListV1Alpha, error) {
	body, status, err := c.BaseClient.List(c.ResourceNamePlural)

	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, fmt.Errorf(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewDashboardListV1AlphaFromJson(body)
}

func (c *DashboardApiV1AlphaApi) GetDashboard(name string) (*models.DashboardV1Alpha, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, name)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewDashboardV1AlphaFromJson(body)
}

func (c *DashboardApiV1AlphaApi) DeleteDashboard(name string) error {
	body, status, err := c.BaseClient.Delete(c.ResourceNamePlural, name)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return nil
}

func (c *DashboardApiV1AlphaApi) CreateDashboard(d *models.DashboardV1Alpha) (*models.DashboardV1Alpha, error) {
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

	return models.NewDashboardV1AlphaFromJson(body)
}

func (c *DashboardApiV1AlphaApi) UpdateDashboard(d *models.DashboardV1Alpha) (*models.DashboardV1Alpha, error) {
	json_body, err := d.ToJson()

	if err != nil {
		return nil, fmt.Errorf("failed to serialize %s object '%s'", c.ResourceNameSingular, err)
	}

	identifier := ""

	if d.Metadata.Id != "" {
		identifier = d.Metadata.Id
	} else {
		identifier = d.Metadata.Name
	}

	body, status, err := c.BaseClient.Patch(c.ResourceNamePlural, identifier, json_body)

	if err != nil {
		return nil, fmt.Errorf("updating %s on Semaphore failed '%s'", c.ResourceNamePlural, err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewDashboardV1AlphaFromJson(body)
}
