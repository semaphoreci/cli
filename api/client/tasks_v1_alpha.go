package client

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	models "github.com/semaphoreci/cli/api/models"
	retry "github.com/semaphoreci/cli/api/retry"
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

// ListTasks fetches all tasks for a project using pagination and aggregates them.
func (c *TasksApiV1AlphaApi) ListTasks(projectID string) (models.TaskListV1Alpha, error) {
	query := url.Values{}
	query.Add("project_id", projectID)
	query.Add("page_size", "200")

	allTasks := make(models.TaskListV1Alpha, 0)
	currentPage := 1
	const maxFailures = 5
	// maxTaskPages caps pagination depth to prevent runaway loops;
	// at 200 items/page this allows up to 100k tasks per project.
	const maxTaskPages = 500

	for {
		query.Set("page", fmt.Sprintf("%d", currentPage))

		var page models.TaskListV1Alpha
		var headers http.Header

		err := retry.RetryWithMaxFailures(maxFailures, func() error {
			page = nil
			headers = nil

			body, status, hdrs, err := c.BaseClient.ListWithParams(c.ResourceNamePlural, query)
			headers = hdrs
			if err != nil {
				return fmt.Errorf("connecting to Semaphore failed: %w", err)
			}
			if status != http.StatusOK {
				msg := string(body)
				if len(msg) > 200 {
					msg = msg[:200] + "...(truncated)"
				}
				httpErr := fmt.Errorf("http status %d with message \"%s\" received from upstream", status, msg)
				if status >= 300 && status < 500 && status != http.StatusTooManyRequests {
					return retry.NonRetryable(httpErr)
				}
				return httpErr
			}

			pageList, err := models.NewTaskListV1AlphaFromJSON(body)
			if err != nil {
				return retry.NonRetryable(fmt.Errorf("failed to deserialize tasks list: %w", err))
			}
			page = pageList
			return nil
		})

		if err != nil {
			return nil, fmt.Errorf("failed fetching page %d (after accumulating %d tasks from %d pages): %w",
				currentPage, len(allTasks), currentPage-1, err)
		}

		if headers == nil {
			return nil, fmt.Errorf("internal error: response headers missing after fetching page %d (accumulated %d tasks)", currentPage, len(allTasks))
		}

		allTasks = append(allTasks, page...)

		if !hasNextPage(headers) {
			break
		}
		if currentPage >= maxTaskPages {
			return nil, fmt.Errorf("pagination safety limit reached (%d pages); results may be incomplete -- please narrow your query", maxTaskPages)
		}
		currentPage++
	}

	return allTasks, nil
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
