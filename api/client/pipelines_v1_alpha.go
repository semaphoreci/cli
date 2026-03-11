package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	models "github.com/semaphoreci/cli/api/models"
	retry "github.com/semaphoreci/cli/api/retry"
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

	var allPipelines models.PipelinesListV1Alpha
	currentPage := 1
	const maxFailures = 5
	const maxPages = 500
	query.Add("page_size", "200")

	for {
		query.Set("page", fmt.Sprintf("%d", currentPage))

		var page models.PipelinesListV1Alpha
		var headers http.Header

		err := retry.RetryWithMaxFailures(maxFailures, func() error {
			body, status, hdrs, err := c.BaseClient.ListWithParams(c.ResourceNamePlural, query)
			headers = hdrs
			if err != nil {
				return fmt.Errorf("connecting to Semaphore failed '%s'", err)
			}
			if status != http.StatusOK {
				return fmt.Errorf("http status %d with message \"%s\" received from upstream", status, string(body))
			}
			if err := page.UnmarshalJSON(body); err != nil {
				return fmt.Errorf("failed to deserialize pipelines list '%s'", err)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		allPipelines = append(allPipelines, page...)

		if !hasNextPage(headers) || currentPage >= maxPages {
			break
		}
		currentPage++
	}

	return allPipelines.MarshalJSON()
}

// hasNextPage checks for a Link header with rel="next" to determine
// if more pages are available.
func hasNextPage(headers http.Header) bool {
	for _, link := range headers.Values("Link") {
		for _, part := range strings.Split(link, ",") {
			if strings.Contains(part, `rel="next"`) {
				return true
			}
		}
	}
	return false
}
