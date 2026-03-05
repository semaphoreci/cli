package client

import (
	"errors"
	"fmt"
	"net/url"

	models "github.com/semaphoreci/cli/api/models"
)

type TasksApiV1AlphaApi struct {
	BaseClient           BaseClient
	ResourceNameSingular string
	ResourceNamePlural   string
}

func NewTasksV1AlphaApi() TasksApiV1AlphaApi {
	baseClient := NewBaseClientFromConfig()
	baseClient.SetApiVersion("v1alpha")

	return TasksApiV1AlphaApi{
		BaseClient:           baseClient,
		ResourceNamePlural:   "tasks",
		ResourceNameSingular: "task",
	}
}

func (c *TasksApiV1AlphaApi) ListTasks(projectID string) (models.TaskListV1Alpha, error) {
	query := url.Values{}
	query.Add("project_id", projectID)

	body, status, _, err := c.BaseClient.ListWithParams(c.ResourceNamePlural, query)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewTaskListV1AlphaFromJSON(body)
}

func (c *TasksApiV1AlphaApi) DescribeTask(id string) (*models.TaskDescribeV1Alpha, error) {
	body, status, err := c.BaseClient.Get(c.ResourceNamePlural, id)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewTaskDescribeV1AlphaFromJSON(body)
}

func (c *TasksApiV1AlphaApi) RunTask(id string, requestBody []byte) (*models.RunTaskResponse, error) {
	body, status, err := c.BaseClient.PostAction(c.ResourceNamePlural, id, "run_now", requestBody)

	if err != nil {
		return nil, errors.New(fmt.Sprintf("connecting to Semaphore failed '%s'", err))
	}

	if status != 200 {
		return nil, errors.New(fmt.Sprintf("http status %d with message \"%s\" received from upstream", status, body))
	}

	return models.NewRunTaskResponseFromJSON(body)
}
