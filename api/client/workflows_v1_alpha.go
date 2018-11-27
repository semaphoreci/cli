package client

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	models "github.com/semaphoreci/cli/api/models"
)

type WorkflowApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewWorkflowV1AlphaApi() WorkflowApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return WorkflowApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "plumber-workflows",
		ResourceNameSingular: "plumber-workflow",
	}
}

func (c *WorkflowApiV1AlphaApi) ListWorkflows(project_id string) (*models.WorkflowListV1Alpha, error) {
	urlEncode := fmt.Sprintf("%s?project_id=%s", c.ResourceNamePlural, project_id)
	body, status, err := c.BaseClient.List(urlEncode)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewWorkflowListV1AlphaFromJson(body)
}

func (c *WorkflowApiV1AlphaApi) CreateSnapshotWf(project_id, label string, archiveContent []byte) ([]byte, error) {
	requestToken, err := uuid.NewUUID()

	if err != nil {
		return nil, fmt.Errorf("uuid creation failed '%s'", err)
	}

	v := url.Values{}
	v.Set("project_id", project_id)
	v.Set("label", label)
	v.Set("request_token", requestToken.String())
	v.Set("snapshot_archive", string(archiveContent))
	requestBody := v.Encode()

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	body, status, err := c.BaseClient.PostHeaders(c.ResourceNamePlural, []byte(requestBody), headers)

	switch {
	case err != nil:
		return nil, fmt.Errorf("connecting to Semaphore failed '%s'", err)
	case status != 200:
		return nil, fmt.Errorf("http status %d with message \"%s\" received from upstream", status, body)
	}

	return body, nil
}

func (c *WorkflowApiV1AlphaApi) StopWf(id string) ([]byte, error) {                                                                                  
	requestToken, err := uuid.NewUUID()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("request token generation failed '%s'", err))
	}

	actionArgs := fmt.Sprintf("%s?%s=%s", "terminate", "request_token", requestToken.String())
	body, status, err := c.BaseClient.PostAction(c.ResourceNamePlural, id, actionArgs, []byte(""))

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return body, nil
}

func (c *WorkflowApiV1AlphaApi) Rebuild(id string) ([]byte, error) {
	requestToken, err := uuid.NewUUID()

        if err != nil {
		return nil, errors.New(fmt.Sprintf("request token generation failed '%s'", err))
        }

        actionArgs := fmt.Sprintf("%s?%s=%s", "reschedule", "request_token", requestToken.String())
        body, status, err := c.BaseClient.PostAction(c.ResourceNamePlural, id, actionArgs, []byte(""))

        if err != nil {
                return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
        }

        if status != 200 {
                return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
        }

        return body, nil
}

