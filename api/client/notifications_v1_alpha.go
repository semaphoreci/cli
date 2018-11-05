package client

import (
	"errors"
	"fmt"

	models "github.com/semaphoreci/cli/api/models"
)

type NotificationsV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewNotificationsV1AlphaApi() NotificationsV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return NotificationsV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "notifications",
		ResourceNameSingular: "notification",
	}
}

func (c *NotificationsV1AlphaApi) ListNotifications() (*models.DashboardListV1Alpha, error) {
	body, status, err := c.BaseClient.List(c.ResourceNamePlural)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewNotificationListV1AlphaFromJson(body)
}

func (c *NotificationsV1AlphaApi) GetNotification(name string) (*models.NotificationV1Alpha, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, name)

	if err != nil {
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewNotificationV1AlphaFromJson(body)
}

func (c *NotificationsV1AlphaApi) DeleteNotification(name string) error {
	body, status, err := c.BaseClient.Delete(c.ResourceNamePlural, name)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return nil
}

func (c *NotificationsV1AlphaApi) CreateNotification(n *models.NotificationV1Alpha) (*models.NotificationV1Alpha, error) {
	json_body, err := n.ToJson()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to serialize object '%s'", err))
	}

	body, status, err := c.BaseClient.Post(c.ResourceNamePlural, json_body)

	if err != nil {
		return nil, fmt.Errorf("creating %s on Semaphore failed '%s'", c.ResourceNameSingular, err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewNotificationV1AlphaFromJson(body)
}

func (c *NotificationsV1AlphaApi) UpdateNotification(n *models.NotificationV1Alpha) (*models.NotificationV1Alpha, error) {
	json_body, err := n.ToJson()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to serialize %s object '%s'", c.ResourceNameSingular, err))
	}

	identifier := ""

	if n.Metadata.Id != "" {
		identifier = n.Metadata.Id
	} else {
		identifier = n.Metadata.Name
	}

	body, status, err := c.BaseClient.Patch(c.ResourceNamePlural, identifier, json_body)

	if err != nil {
		return nil, fmt.Errorf("updating %s on Semaphore failed '%s'", c.ResourceNamePlural, err)
	}

	if status != 200 {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return models.NewNotificationV1AlphaFromJson(body)
}
