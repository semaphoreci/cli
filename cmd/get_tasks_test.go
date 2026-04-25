package cmd

import (
	"net/http"
	"regexp"
	"testing"

	client "github.com/semaphoreci/cli/api/client"
	httpmock "github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test__ListTasks__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/projects/test-project",
		func(req *http.Request) (*http.Response, error) {
			p := `{
				"metadata": {
					"id": "758cb945-7495-4e40-a9a1-4b3991c6a8fe"
				}
			}`
			return httpmock.NewStringResponse(200, p), nil
		},
	)

	httpmock.RegisterRegexpResponder("GET",
		regexp.MustCompile(`https://org\.semaphoretext\.xyz/api/v1alpha/tasks\?.*project_id=758cb945-7495-4e40-a9a1-4b3991c6a8fe`),
		func(req *http.Request) (*http.Response, error) {
			received = true

			tasks := `[
				{
					"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
					"name": "deploy",
					"project_id": "758cb945-7495-4e40-a9a1-4b3991c6a8fe",
					"branch": "main",
					"pipeline_file": ".semaphore/deploy.yml",
					"recurring": false,
					"parameters": [
						{"name": "ENV", "required": true, "default_value": "staging"}
					]
				},
				{
					"id": "cc3ba294-d4b3-48bc-90a7-12dd56e9424d",
					"name": "nightly",
					"project_id": "758cb945-7495-4e40-a9a1-4b3991c6a8fe",
					"branch": "main",
					"pipeline_file": ".semaphore/nightly.yml",
					"recurring": true,
					"paused": true
				}
			]`

			return httpmock.NewStringResponse(200, tasks), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "tasks", "--project-name", "test-project"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET /tasks?project_id=...")
}

func Test__ListTasks__WithProjectID(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterRegexpResponder("GET",
		regexp.MustCompile(`https://org\.semaphoretext\.xyz/api/v1alpha/tasks\?.*project_id=758cb945-7495-4e40-a9a1-4b3991c6a8fe`),
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, "[]"), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "tasks", "--project-id", "758cb945-7495-4e40-a9a1-4b3991c6a8fe"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET /tasks with project_id flag")
}

func Test__ListTasks__SuspendedTask(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterRegexpResponder("GET",
		regexp.MustCompile(`https://org\.semaphoretext\.xyz/api/v1alpha/tasks\?.*project_id=758cb945-7495-4e40-a9a1-4b3991c6a8fe`),
		func(req *http.Request) (*http.Response, error) {
			received = true

			tasks := `[
				{
					"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
					"name": "deploy",
					"project_id": "758cb945-7495-4e40-a9a1-4b3991c6a8fe",
					"branch": "main",
					"pipeline_file": ".semaphore/deploy.yml",
					"recurring": true,
					"paused": false,
					"suspended": true
				}
			]`

			return httpmock.NewStringResponse(200, tasks), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "tasks", "--project-id", "758cb945-7495-4e40-a9a1-4b3991c6a8fe"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET /tasks with suspended task")
}

func Test__DescribeTask__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
		func(req *http.Request) (*http.Response, error) {
			received = true

			task := `{
				"schedule": {
					"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
					"name": "deploy",
					"project_id": "aa1ba294-d4b3-48bc-90a7-12dd56e9424a",
					"branch": "main",
					"at": "0 0 * * *",
					"pipeline_file": ".semaphore/deploy.yml",
					"recurring": true
				},
				"triggers": [
					{
						"triggered_at": "2024-01-15 09:00:00",
						"scheduling_status": "passed",
						"scheduled_workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e",
						"branch": "main"
					}
				]
			}`

			return httpmock.NewStringResponse(200, task), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "tasks", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET /tasks/:id")
}

func Test__DescribeTask__NoTriggers(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
		func(req *http.Request) (*http.Response, error) {
			received = true

			task := `{
				"schedule": {
					"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
					"name": "deploy",
					"project_id": "aa1ba294-d4b3-48bc-90a7-12dd56e9424a",
					"pipeline_file": ".semaphore/deploy.yml"
				}
			}`

			return httpmock.NewStringResponse(200, task), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "tasks", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET /tasks/:id with no triggers")
}

func Test__GetTasks__TaskAlias(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
		func(req *http.Request) (*http.Response, error) {
			received = true

			task := `{
				"schedule": {
					"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
					"name": "deploy",
					"project_id": "aa1ba294-d4b3-48bc-90a7-12dd56e9424a",
					"pipeline_file": ".semaphore/deploy.yml"
				}
			}`

			return httpmock.NewStringResponse(200, task), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "task", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the 'task' alias to work for describe")
}

func Test__ListTasks__MultiPage(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	page1Received := false
	page2Received := false

	httpmock.RegisterRegexpResponder("GET",
		regexp.MustCompile(`https://org\.semaphoretext\.xyz/api/v1alpha/tasks.*[?&]page=1(?:&|$)`),
		func(req *http.Request) (*http.Response, error) {
			page1Received = true
			body := `[
                {"id": "11111111-1111-1111-1111-111111111111", "name": "t1", "project_id": "758cb945-7495-4e40-a9a1-4b3991c6a8fe", "branch": "main", "pipeline_file": ".semaphore/t1.yml", "recurring": false}
            ]`
			resp := httpmock.NewStringResponse(200, body)
			resp.Header.Set("Link", `<https://org.semaphoretext.xyz/api/v1alpha/tasks?page=2>; rel="next"`)
			return resp, nil
		},
	)

	httpmock.RegisterRegexpResponder("GET",
		regexp.MustCompile(`https://org\.semaphoretext\.xyz/api/v1alpha/tasks.*[?&]page=2(?:&|$)`),
		func(req *http.Request) (*http.Response, error) {
			page2Received = true
			body := `[
                {"id": "22222222-2222-2222-2222-222222222222", "name": "t2", "project_id": "758cb945-7495-4e40-a9a1-4b3991c6a8fe", "branch": "main", "pipeline_file": ".semaphore/t2.yml", "recurring": false}
            ]`
			return httpmock.NewStringResponse(200, body), nil
		},
	)

	c := client.NewTasksV1AlphaApi()
	tasks, err := c.ListTasks("758cb945-7495-4e40-a9a1-4b3991c6a8fe")

	assert.NoError(t, err)
	assert.True(t, page1Received, "Expected page 1 to be fetched")
	assert.True(t, page2Received, "Expected page 2 to be fetched (pagination must follow Link: rel=next)")
	assert.Len(t, tasks, 2, "Expected tasks from both pages to be aggregated")
	assert.Equal(t, "t1", tasks[0].Name, "Expected first task from page 1")
	assert.Equal(t, "t2", tasks[1].Name, "Expected second task from page 2")
	assert.Equal(t, 2, httpmock.GetTotalCallCount(), "Expected exactly two HTTP requests; pagination must stop when Link: rel=next is absent")
}
