package client

import (
	"errors"
	"fmt"
	"net/url"

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

func (c *JobsApiV1AlphaApi) ListJobs(states []string) (*models.JobListV1Alpha, error) {
	query := url.Values{}

	for _, s := range states {
		query.Add("states", s)
	}

	body, status, err := c.BaseClient.ListWithParams(c.ResourceNamePlural, query)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewJobListV1AlphaFromJson(body)
}

func (c *JobsApiV1AlphaApi) GetJob(name string) (*models.JobV1Alpha, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, name)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewJobV1AlphaFromJson(body)
}

func (c *JobsApiV1AlphaApi) GetJobDebugSSHKey(id string) (*models.JobDebugSSHKeyV1Alpha, error) {
	path := fmt.Sprintf("%s/%s", id, "debug_ssh_key")
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, path)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewJobDebugSSHKeyV1AlphaFromJSON(body)
}

func (c *JobsApiV1AlphaApi) CreateJob(j *models.JobV1Alpha) (*models.JobV1Alpha, error) {
	json_body, err := j.ToJson()

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

	return models.NewJobV1AlphaFromJson(body)
}

func (c *JobsApiV1AlphaApi) StopJob(id string) error {
	path := fmt.Sprintf("%s/%s/%s", c.ResourceNamePlural, id, "stop")
	body, status, err := c.BaseClient.Post(path, []byte{})

	if err != nil {
		return errors.New(fmt.Sprintf("stopping %s on Semaphore failed '%s'", c.ResourceNameSingular, err))
	}

	if status != 200 {
		return errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return nil
}
