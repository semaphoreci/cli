package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/api/uuid"
	assert "github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func Test__EditDashboard__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dash := `{
		"metadata": {
			"name":"my-work",
			"id":"07e64c23-4bca-4f62-b76e-4a85aef24bbb",
			"title":"My Work",
			"create_time":"1536843668",
			"update_time":"1536843668"
		},
		"spec": {
			"widgets": [{
				"name": "Workflows",
				"type": "list",
				"filters":{
					"github_uid":"{{github_uid}}"
				}
			}]
		}
	}`

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/dashboards/my-work",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	httpmock.RegisterResponder("PATCH", "https://org.semaphoretext.xyz/api/v1alpha/dashboards/07e64c23-4bca-4f62-b76e-4a85aef24bbb",
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "dashboards", "my-work"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET and PATCH dashboard")
	}
}

func Test__EditNotification__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var received *models.NotificationV1Alpha

	endpoint := "https://org.semaphoretext.xyz/api/v1alpha/notifications/test"

	httpmock.RegisterResponder("GET", endpoint,
		func(req *http.Request) (*http.Response, error) {
			data, _ := ioutil.ReadFile("../fixtures/notification.yml")
			notif, _ := models.NewNotificationV1AlphaFromYaml(data)
			json, _ := notif.ToJson()

			return httpmock.NewStringResponse(200, string(json)), nil
		},
	)

	httpmock.RegisterResponder("PATCH", endpoint,
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			received, _ = models.NewNotificationV1AlphaFromJson(body)

			return httpmock.NewStringResponse(200, string(body)), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "notification", "test"})
	RootCmd.Execute()

	assert.Equal(t, received.Metadata.Name, "test")

	rule := received.Spec.Rules[0]

	assert.Equal(t, rule.Name, "Rule #1")
	assert.Equal(t, rule.Filter.Projects, []string{"cli"})
	assert.Equal(t, rule.Filter.Branches, []string{"master"})
	assert.Equal(t, rule.Filter.Pipelines, []string{"semaphore.yml"})
}

func Test__EditSecret__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dash := `{
		"metadata":{
			"name":"aaaaaaa",
			"id":"bb2ba294-d4b3-48bc-90a7-12dd56e9424b",
			"create_time":"1536673464",
			"update_time":"1536674946"
		},
		"data":{
			"env_vars":[{
				"name":"TEST",
				"value":"AAAA"
			}],
			"files":[{
				"path":"a.txt",
				"content":"W1twYWNrYWdlXV0KbmFtZSA9ICJzcGlyYWwtbWVtb3J5Igp2ZXJzaW9uID0gIjAuMS4wIgoK"
			}]
		}
	}`

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1beta/secrets/aaaaaaa",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	httpmock.RegisterResponder("PATCH", "https://org.semaphoretext.xyz/api/v1beta/secrets/bb2ba294-d4b3-48bc-90a7-12dd56e9424b",
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "secrets", "aaaaaaa"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET and PATCH secret")
	}
}

func Test__EditProject__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dash := `{
		"metadata":{
			"name":"hello",
			"id":"bb2ba294-d4b3-48bc-90a7-12dd56e9424b",
			"description":"Just saying hi!"
		},
		"spec":{
			"repository":{
				"url":"git@github.com/renderextext/hello",
				"run_on":["tags", "branches"],
				"forked_pull_requests":{
					"allowed_secrets":["foo"]
				},
				"pipeline_file": ""
			},
			"schedulers":[
				{
					"name":"cron",
					"id":"bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
					"branch":"master",
					"at":"* * * *",
					"pipeline_file":".semaphore/cron.yml"
				}
			]
		}
	}`

	var received *models.ProjectV1Alpha

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/projects/hello",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	httpmock.RegisterResponder("PATCH", "https://org.semaphoretext.xyz/api/v1alpha/projects/bb2ba294-d4b3-48bc-90a7-12dd56e9424b",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			received, _ = models.NewProjectV1AlphaFromJson(body)

			return httpmock.NewStringResponse(200, string(body)), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "project", "hello"})
	RootCmd.Execute()

	assert.Equal(t, received.Metadata.Name, "hello")
	assert.Equal(t, received.Metadata.Description, "Just saying hi!")

	repo := received.Spec.Repository

	assert.Equal(t, repo.Url, "git@github.com/renderextext/hello")
	assert.Equal(t, repo.RunOn, []string{"tags", "branches"})

	forked_pull_requests := received.Spec.Repository.ForkedPullRequests

	assert.Equal(t, forked_pull_requests.AllowedSecrets, []string{"foo"})

	scheduler := received.Spec.Schedulers[0]

	assert.Equal(t, scheduler.Name, "cron")
	assert.Equal(t, scheduler.Branch, "master")
	assert.Equal(t, scheduler.At, "* * * *")
	assert.Equal(t, scheduler.PipelineFile, ".semaphore/cron.yml")
}

func Test__EditProject__WithTasks__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	dash := `{
		"metadata":{
			"name":"hello",
			"id":"bb2ba294-d4b3-48bc-90a7-12dd56e9424b",
			"description":"Just saying hi!"
		},
		"spec":{
			"repository":{
				"url":"git@github.com/renderextext/hello",
				"run_on":["tags", "branches"],
				"forked_pull_requests":{
					"allowed_secrets":["foo"]
				},
				"pipeline_file": ""
			},
			"tasks":[
				{
					"name":"cron",
					"description":"cron description",
					"id":"bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
					"scheduled":false,
					"branch":"master",
					"pipeline_file":".semaphore/cron.yml",
					"parameters":[
						{
							"name":"param1",
							"required":true,
							"description":"param1 description",
							"default_value":"option1",
							"options":["option1", "option2"]
						}
					]
				}
			]
		}
	}`

	var received *models.ProjectV1Alpha

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/projects/hello",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	httpmock.RegisterResponder("PATCH", "https://org.semaphoretext.xyz/api/v1alpha/projects/bb2ba294-d4b3-48bc-90a7-12dd56e9424b",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			received, _ = models.NewProjectV1AlphaFromJson(body)

			return httpmock.NewStringResponse(200, string(body)), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "project", "hello"})
	RootCmd.Execute()

	assert.Equal(t, received.Metadata.Name, "hello")
	assert.Equal(t, received.Metadata.Description, "Just saying hi!")

	repo := received.Spec.Repository

	assert.Equal(t, repo.Url, "git@github.com/renderextext/hello")
	assert.Equal(t, repo.RunOn, []string{"tags", "branches"})

	forked_pull_requests := received.Spec.Repository.ForkedPullRequests

	assert.Equal(t, forked_pull_requests.AllowedSecrets, []string{"foo"})

	task := received.Spec.Tasks[0]

	assert.Equal(t, task.Name, "cron")
	assert.Equal(t, task.Description, "cron description")
	assert.Equal(t, task.Branch, "master")
	assert.Equal(t, task.Scheduled, false)
	assert.Equal(t, task.PipelineFile, ".semaphore/cron.yml")

	task_parameter := task.Parameters[0]

	assert.Equal(t, task_parameter.Name, "param1")
	assert.Equal(t, task_parameter.Required, true)
	assert.Equal(t, task_parameter.Description, "param1 description")
	assert.Equal(t, task_parameter.DefaultValue, "option1")
	assert.Equal(t, task_parameter.Options, []string{"option1", "option2"})
}

func Test__EditDeploymentTarget__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	uuid.Mock()
	defer uuid.Unmock()

	targetJSON := `{
			"apiVersion": "v1alpha",
			"kind": "DeploymentTarget",
			"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
			"name": "dt-name",
			"project_id": "projId1",
			"organization_id": "org-id",
			"description": "dt-description",
			"url": "www.semaphore.xyz",
			"subject_rules": [
			{
				"type": "USER",
				"subject_id": "00000000-0000-0000-0000-000000000000"
			}
			],
			"object_rules": [
			{
				"type": "BRANCH",
				"match_mode": "PATTERN",
				"pattern": ".*main.*"
			}
			],
			"env_vars": [
			{
				"name": "X",
				"value": "123"
			}
			],
			"active": true,
			"bookmark_parameter1": "book1"
		}
		`

	var received *models.DeploymentTargetV1Alpha

	targetId := "bb2ba294-d4b3-48bc-90a7-12dd56e9424c"
	targetGetURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets/%s?include_secrets=true", targetId)
	httpmock.RegisterResponder(http.MethodGet, targetGetURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, targetJSON), nil
		},
	)
	targetPatchURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets/%s", targetId)
	httpmock.RegisterResponder(http.MethodPatch, targetPatchURL,
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			updateRequest := models.DeploymentTargetCreateRequestV1Alpha{}
			json.Unmarshal(body, &updateRequest)
			received = &updateRequest.DeploymentTargetV1Alpha
			body, _ = json.Marshal(received)
			return httpmock.NewStringResponse(200, string(body)), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "dt", targetId})
	RootCmd.Execute()

	assert.Equal(t, received.Name, "dt-name")
	assert.Equal(t, received.Description, "dt-description")
	assert.Equal(t, received.Url, "www.semaphore.xyz")
	assert.Equal(t, len(*received.EnvVars), 1)
	assert.Equal(t, *(*received.EnvVars)[0], models.DeploymentTargetEnvVarV1Alpha{Name: "X", Value: "123"})
}

func Test__EditDeploymentTargetByName__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	uuid.Mock()
	defer uuid.Unmock()

	targetJSON := `{
			"apiVersion": "v1alpha",
			"kind": "DeploymentTarget",
			"id": "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
			"name": "dt-name",
			"project_id": "projId1",
			"organization_id": "org-id",
			"description": "dt-description",
			"url": "www.semaphore.xyz",
			"subject_rules": [
			{
				"type": "USER",
				"subject_id": "00000000-0000-0000-0000-000000000000"
			}
			],
			"object_rules": [
			{
				"type": "BRANCH",
				"match_mode": "PATTERN",
				"pattern": ".*main.*"
			}
			],
			"env_vars": [
			{
				"name": "X",
				"value": "123"
			}
			],
			"active": true,
			"bookmark_parameter1": "book1"
		}
		`
	targetsJSON := "[" + targetJSON + "]"
	var received *models.DeploymentTargetV1Alpha

	targetId := "bb2ba294-d4b3-48bc-90a7-12dd56e9424c"
	targetName := "dt-name"
	projectId := "projId1"
	targetGetURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets?project_id=%s&target_name=%s", projectId, targetName)
	httpmock.RegisterResponder(http.MethodGet, targetGetURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, targetsJSON), nil
		},
	)
	targetGetByIDURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets/%s?include_secrets=true", targetId)
	httpmock.RegisterResponder(http.MethodGet, targetGetByIDURL,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, targetJSON), nil
		},
	)

	targetPatchURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets/%s", targetId)
	httpmock.RegisterResponder(http.MethodPatch, targetPatchURL,
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			updateRequest := models.DeploymentTargetCreateRequestV1Alpha{}
			json.Unmarshal(body, &updateRequest)
			received = &updateRequest.DeploymentTargetV1Alpha

			body, _ = json.Marshal(received)
			return httpmock.NewStringResponse(200, string(body)), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "dt", targetName, "-i", projectId})
	RootCmd.Execute()

	assert.Equal(t, received.Name, "dt-name")
	assert.Equal(t, received.Description, "dt-description")
	assert.Equal(t, received.Url, "www.semaphore.xyz")
	assert.Equal(t, len(*received.EnvVars), 1)
	assert.Equal(t, *(*received.EnvVars)[0], models.DeploymentTargetEnvVarV1Alpha{Name: "X", Value: "123"})
}

func Test__EditDeploymentTargetDeactivate__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	targetId := "494b76aa-f3f0-4ecf-b5ef-c389591a01be"
	patchURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets/%s/deactivate", targetId)
	httpmock.RegisterResponder(http.MethodPatch, patchURL,
		func(req *http.Request) (*http.Response, error) {
			received = true

			target := `{
				"target_id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
				"cordoned": true
				}
				`

			return httpmock.NewStringResponse(200, target), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "dt", targetId, "-d"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive PATCH deployment_targets/:id/deactivate")
}

func Test__EditDeploymentTargetActivate__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	targetId := "494b76aa-f3f0-4ecf-b5ef-c389591a01be"
	patchURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets/%s/activate", targetId)
	httpmock.RegisterResponder(http.MethodPatch, patchURL,
		func(req *http.Request) (*http.Response, error) {
			received = true

			target := `{
				"target_id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
				"active": false
				}
				`

			return httpmock.NewStringResponse(200, target), nil
		},
	)

	RootCmd.SetArgs([]string{"edit", "dt", targetId, "--activate"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive PATCH deployment_targets/:id/activate")
}
