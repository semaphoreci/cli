package client

import (
	"errors"
	"fmt"
	"net/url"

	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/api/uuid"
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
	detailed := fmt.Sprintf("%s?detailed=true", id)
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, detailed)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewPipelineV1AlphaFromJson(body)
}

func (c *PipelinesApiV1AlphaApi) StopPpl(id string) ([]byte, error) {
	request_body := []byte("{\"terminate_request\": true}")

	body, status, err := c.BaseClient.Patch(c.ResourceNamePlural, id, request_body)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return body, nil
}

func (c *PipelinesApiV1AlphaApi) PartialRebuildPpl(id string) ([]byte, error) {
	requestToken, err := uuid.NewUUID()

	if err != nil {
		return nil, errors.New(fmt.Sprintf("request token generation failed '%s'", err))
	}

	actionArgs := fmt.Sprintf("%s?%s=%s", "partial_rebuild", "request_token", requestToken.String())
	body, status, err := c.BaseClient.PostAction(c.ResourceNamePlural, id, actionArgs, []byte(""))

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return body, nil
}

func (c *PipelinesApiV1AlphaApi) ListPplByWfID(projectID, wfID string) ([]byte, error) {
	detailed := fmt.Sprintf("%s?project_id=%s&wf_id=%s", c.ResourceNamePlural, projectID, wfID)
	body, status, err := c.BaseClient.List(detailed)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return body, nil
}

type ListOptions struct {
	CreatedAfter  int64
	CreatedBefore int64
}

func (c *PipelinesApiV1AlphaApi) ListPpl(projectID string) ([]byte, error) {
	detailed := fmt.Sprintf("%s?project_id=%s", c.ResourceNamePlural, projectID)
	body, status, err := c.BaseClient.List(detailed)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return body, nil
}

func (c *PipelinesApiV1AlphaApi) ListPplWithOptions(projectID string, options ListOptions) ([]byte, error) {
	query := url.Values{}
	query.Add("project_id", projectID)

	if options.CreatedAfter > 0 {
		query.Add("created_after", fmt.Sprintf("%d", options.CreatedAfter))
	}

	if options.CreatedBefore > 0 {
		query.Add("created_before", fmt.Sprintf("%d", options.CreatedBefore))
	}

	body, status, _, err := c.BaseClient.ListWithParams(c.ResourceNamePlural, query)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return body, nil
}
