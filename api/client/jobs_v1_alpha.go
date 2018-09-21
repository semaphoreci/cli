package client

import (
	"errors"
	"fmt"

	models "github.com/semaphoreci/cli/api/models"
)

type JobsApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewJobsV1AlphaApi() JobsApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return JobsApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "jobs",
		ResourceNameSingular: "job",
	}
}

func (c *JobsApiV1AlphaApi) ListJobs() (*models.JobsListV1Alpha, error) {
	body, status, err := c.BaseClient.List(c.ResourceNamePlural)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewJobsListV1AlphaFromJson(body)
}

func (c *JobsApiV1AlphaApi) GetJob(name string) (*models.JobsV1Alpha, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, name)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewJobsV1AlphaFromJson(body)
}
