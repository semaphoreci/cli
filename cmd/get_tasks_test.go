package cmd

import (
	"net/http"
	"testing"

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

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/tasks?project_id=758cb945-7495-4e40-a9a1-4b3991c6a8fe",
		func(req *http.Request) (*http.Response, error) {
			received = true

			tasks := `[
				{
					"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
					"name": "deploy",
					"project_id": "758cb945-7495-4e40-a9a1-4b3991c6a8fe",
					"branch": "main",
					"pipeline_file": ".semaphore/deploy.yml",
					"recurring": false
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

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/tasks?project_id=758cb945-7495-4e40-a9a1-4b3991c6a8fe",
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, "[]"), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "tasks", "--project-id", "758cb945-7495-4e40-a9a1-4b3991c6a8fe"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET /tasks with project_id flag")
}

func Test__ListTasks__EmptyList(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/tasks?project_id=758cb945-7495-4e40-a9a1-4b3991c6a8fe",
		func(req *http.Request) (*http.Response, error) {
			received = true
			return httpmock.NewStringResponse(200, "[]"), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "tasks", "--project-id", "758cb945-7495-4e40-a9a1-4b3991c6a8fe"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET /tasks for empty list")
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
