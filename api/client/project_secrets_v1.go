package client

import (
	"errors"
	"fmt"
	"log"

	models "github.com/semaphoreci/cli/api/models"
)

type ProjectSecretsApiV1Api struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewProjectSecretV1Api(projectID string) ProjectSecretsApiV1Api {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1")

	return ProjectSecretsApiV1Api{
		BaseClient:           baseClient,
		ResourceNamePlural:   fmt.Sprintf("projects/%s/secrets", projectID),
		ResourceNameSingular: fmt.Sprintf("projects/%s/secret", projectID),
	}
}

func (c *ProjectSecretsApiV1Api) ListSecrets() (*models.SecretListV1Beta, error) {
	body, status, err := c.BaseClient.List(c.ResourceNamePlural)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewSecretListV1BetaFromJson(body)
}

func (c *ProjectSecretsApiV1Api) GetSecret(name string) (*models.ProjectSecretV1, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, name)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewProjectSecretV1FromJson(body)
}

func (c *ProjectSecretsApiV1Api) DeleteSecret(name string) error {
	body, status, err := c.BaseClient.Delete(c.ResourceNamePlural, name)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return nil
}

func (c *ProjectSecretsApiV1Api) CreateSecret(d *models.ProjectSecretV1) (*models.ProjectSecretV1, error) {
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

	return models.NewProjectSecretV1FromJson(body)
}

func (c *ProjectSecretsApiV1Api) UpdateSecret(d *models.ProjectSecretV1) (*models.ProjectSecretV1, error) {
	json_body, err := d.ToJson()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to serialize %s object '%s'", c.ResourceNameSingular, err))
	}

	identifier := ""

	if d.Metadata.Id != "" {
		identifier = d.Metadata.Id
	} else {
		identifier = d.Metadata.Name
	}

	body, status, err := c.BaseClient.Patch(c.ResourceNamePlural, identifier, json_body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("updating %s on Semaphore failed '%s'", c.ResourceNamePlural, err))
	}

	if status != 200 {
		fallbackResponse, err := c.fallbackUpdate(identifier, d)
		if err != nil {
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
		}
		return fallbackResponse, nil
	}

	return models.NewProjectSecretV1FromJson(body)
}

func (c *ProjectSecretsApiV1Api) fallbackUpdate(identifier string, d *models.ProjectSecretV1) (*models.ProjectSecretV1, error) {
	err := c.DeleteSecret(identifier)

	if err != nil {
		log.Println("fallbackUpdate:", err)
		return nil, fmt.Errorf("updating %s on Semaphore failed '%s'", c.ResourceNamePlural, err)
	}

	return c.CreateSecret(d)
}