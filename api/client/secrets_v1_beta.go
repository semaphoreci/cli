package client

import (
	"errors"
	"fmt"
	"log"

	models "github.com/semaphoreci/cli/api/models"
)

type SecretApiV1BetaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewSecretV1BetaApi() SecretApiV1BetaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1beta")

	return SecretApiV1BetaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "secrets",
		ResourceNameSingular: "secret",
	}
}

func (c *SecretApiV1BetaApi) ListSecrets() (*models.SecretListV1Beta, error) {
	body, status, err := c.BaseClient.List(c.ResourceNamePlural)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewSecretListV1BetaFromJson(body)
}

func (c *SecretApiV1BetaApi) GetSecret(name string) (*models.SecretV1Beta, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, name)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewSecretV1BetaFromJson(body)
}

func (c *SecretApiV1BetaApi) DeleteSecret(name string) error {
	body, status, err := c.BaseClient.Delete(c.ResourceNamePlural, name)

	if err != nil {
		return err
	}

	if status != 200 {
		return fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return nil
}

func (c *SecretApiV1BetaApi) CreateSecret(d *models.SecretV1Beta) (*models.SecretV1Beta, error) {
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

	return models.NewSecretV1BetaFromJson(body)
}

func (c *SecretApiV1BetaApi) UpdateSecret(d *models.SecretV1Beta) (*models.SecretV1Beta, error) {
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

	fmt.Printf("This might overwrite the secret '%s' and change it's id.", identifier)

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

	return models.NewSecretV1BetaFromJson(body)
}

func (c *SecretApiV1BetaApi) fallbackUpdate(identifier string, d *models.SecretV1Beta) (*models.SecretV1Beta, error) {
	err := c.DeleteSecret(identifier)

	if err != nil {
		log.Println("fallbackUpdate:", err)
		return nil, fmt.Errorf("updating %s on Semaphore failed '%s'", c.ResourceNamePlural, err)
	}

	return c.CreateSecret(d)
}