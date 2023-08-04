package cmd

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test__ListProjects__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/projects",
		func(req *http.Request) (*http.Response, error) {
			received = true

			p1 := `{
				"spec": {
					"repository": {
						"url" : "git@github.com:shiroyasha/advent-of-code-2017.git"
					}
				},
				"metadata": {
					"name":"advent-of-code-2017",
					"id":"8f100520-5ab9-469f-854a-87bae95f19b9"
				},
				"kind":"Project",
				"apiVersion":"v1alpha"
			}`

			p2 := `{
				"spec": {
					"repository": {
						"url" : "git@github.com:shiroyasha/test.git"
					}
				},
				"metadata": {
					"name":"test",
					"id":"1f100520-5ab9-469f-854a-87bae95f19b9"
				},
				"kind":"Project",
				"apiVersion":"v1alpha"
			}`

			projects := fmt.Sprintf("[%s,%s]", p1, p2)

			return httpmock.NewStringResponse(200, projects), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "projects"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET projects")
	}
}

func Test__GetProject__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/projects/advent",
		func(req *http.Request) (*http.Response, error) {
			received = true

			p1 := `{
				"spec": {
					"repository": {
						"url" : "git@github.com:shiroyasha/advent-of-code-2017.git"
					}
				},
				"metadata": {
					"name":"advent",
					"id":"8f100520-5ab9-469f-854a-87bae95f19b9"
				},
				"kind":"Project",
				"apiVersion":"v1alpha"
			}`

			return httpmock.NewStringResponse(200, p1), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "projects", "advent"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET projects/test-prj")
	}
}

func Test__ListDashboards__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/dashboards",
		func(req *http.Request) (*http.Response, error) {
			received = true

			d1 := `{
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

			d2 := `{
				"metadata":{
					"name":"everyones-activity",
					"id":"5aabd382-a2b6-46d1-abb5-06d383dace08",
					"title":"Everyone’s Activity",
					"create_time":"1536843668",
					"update_time":"1536843668"
				},
				"spec":{
					"widgets":[{
						"name":"Workflows",
						"type":"list"
					}]
				}
			}`

			dashboards := fmt.Sprintf(`{"dashboards":[%s,%s]}`, d1, d2)

			return httpmock.NewStringResponse(200, dashboards), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "dashboards"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET dashboards")
	}
}

func Test__GetDashboard__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/dashboards/my-work",
		func(req *http.Request) (*http.Response, error) {
			received = true

			d1 := `{
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

			return httpmock.NewStringResponse(200, d1), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "dashboards", "my-work"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET dashboards/my-work")
	}
}

func Test__ListSecrets__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1beta/secrets",
		func(req *http.Request) (*http.Response, error) {
			received = true

			s1 := `{
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

			s2 := `{
				"metadata":{
					"name":"aaa-bbb",
					"id":"d4cdb2aa-e9c6-4077-bf27-721c3a7993a9",
					"create_time":"1536673704",
					"update_time":"1536673704"
				},
				"data":{
				}
			}`

			secrets := fmt.Sprintf(`{"secrets":[%s,%s]}`, s1, s2)

			return httpmock.NewStringResponse(200, secrets), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "secrets"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET secrets")
	}
}

func Test__GetSecret__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1beta/secrets/aaaaaaa",
		func(req *http.Request) (*http.Response, error) {
			received = true

			s1 := `{
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

			return httpmock.NewStringResponse(200, s1), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "secrets", "aaaaaaa"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET secrets/aaaaaaa")
	}
}

func Test__GetAgentType__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/self_hosted_agent_types/s1-testing",
		func(req *http.Request) (*http.Response, error) {
			received = true

			a1 := `{
				"metadata":{
					"name":"s1-testing",
					"create_time":1536673464,
					"update_time":1536674946
				},
				"status":{
					"total_agent_count":0
				}
			}`

			return httpmock.NewStringResponse(200, a1), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "agent_type", "s1-testing"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET self_hosted_agent_types/s1-testing")
	}
}

func Test__GetAgent__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/agents/asdfasdfadsf",
		func(req *http.Request) (*http.Response, error) {
			received = true

			a1 := `{
				"metadata":{
					"name":"s1-testing",
					"type":"asdfasdfadsf",
					"connected_at":1536673464,
					"disabled_at":1536674946
				},
				"status":{
					"state":"waiting_for_job"
				}
			}`

			return httpmock.NewStringResponse(200, a1), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "agent", "asdfasdfadsf"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET agents/asdfasdfadsf")
	}
}

func Test__GetPipeline__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/pipelines/494b76aa-f3f0-4ecf-b5ef-c389591a01be",
		func(req *http.Request) (*http.Response, error) {
			received = true

			p := `{
				"pipeline": {
					"ppl_id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
					"name": "snapshot test",
				"state": "done",
				"result": "passed",
					"result_reason": "test",
				"error_description": ""
				}
			}`

			return httpmock.NewStringResponse(200, p), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "pipelines", "494b76aa-f3f0-4ecf-b5ef-c389591a01be"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET pipelines/:id")
}

func Test__GetDeploymentTargetByName__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	projectId := "494b76aa-f3f0-4ecf-b5ef-c389591a01be"
	targetName := "dep target test"
	getURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets?project_id=%s&target_name=%s", projectId, url.PathEscape(targetName))

	httpmock.RegisterResponder(http.MethodGet, getURL,
		func(req *http.Request) (*http.Response, error) {
			received = true

			p := `[{
						"id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
					"name": "dep target test",
					"url": "https://semaphoreci.xyz/target",
					"project_id": "proj_id"
			}]
			`

			return httpmock.NewStringResponse(200, p), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "dt", targetName, "-i", projectId})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET deployment_targets")
}

func Test__GetDeploymentTarget__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	targetId := "494b76aa-f3f0-4ecf-b5ef-c389591a01be"
	getURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets/%s", targetId)
	httpmock.RegisterResponder(http.MethodGet, getURL,
		func(req *http.Request) (*http.Response, error) {
			received = true

			p := `{
						"id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
					"name": "dep target test",
					"url": "https://semaphoreci.xyz/target",
					"project_id": "proj_id"
			}
			`

			return httpmock.NewStringResponse(200, p), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "dt", targetId})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET deployment_targets/:id")
}

func Test__GetDeploymentTargetsList__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	projectId := "proj_id"
	getURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets?project_id=%s", projectId)
	httpmock.RegisterResponder(http.MethodGet, getURL,
		func(req *http.Request) (*http.Response, error) {
			received = true

			p := `[{
						"id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
					"name": "dep target test",
					"url": "https://semaphoreci.xyz/target",
					"project_id": "proj_id"
			}]
			`

			return httpmock.NewStringResponse(200, p), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "dt", "--project-id", projectId})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET deployment_targets/:id")
}

func Test__GetDeploymentTargetHistory__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	targetId := "494b76aa-f3f0-4ecf-b5ef-c389591a01be"
	getURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets/%s/history?cursor_type=AFTER&cursor_value=123123123123&git_ref_label=main&git_ref_type=branch&triggered_by=event", targetId)
	httpmock.RegisterResponder(http.MethodGet, getURL,
		func(req *http.Request) (*http.Response, error) {
			received = true

			p := `{"deployments":[{
						"id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
					"target_id": "target_id123",
					"pipeline_id": "pipeline_id"
			}]}`
			return httpmock.NewStringResponse(200, p), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "dt", targetId,
		"--history",
		"-a", "123123123123",
		"--triggered-by", "event",
		"--git-ref-type", "branch",
		"--git-ref-label", "main",
		"-p", "bookmark 1",
	})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET deployment_targets/:id/history")
}

func Test__GetWorkflows__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/projects/foo",
		func(req *http.Request) (*http.Response, error) {
			received = true

			p := `{
				"metadata": {
					"id": "758cb945-7495-4e40-a9a1-4b3991c6a8fe"
				}
			}`

			return httpmock.NewStringResponse(200, p), nil
		},
	)

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/plumber-workflows?project_id=758cb945-7495-4e40-a9a1-4b3991c6a8fe",
		func(req *http.Request) (*http.Response, error) {
			received = true

			p := `[{
				"wf_id": "b129e277-4aa5-4308-8e31-ec825815e335",
				"requester_id": "92f81b82-3584-4852-ab28-4866624bed1e",
				"project_id": "758cb945-7495-4e40-a9a1-4b3991c6a8fe",
				"initial_ppl_id": "92f81b82-3584-4852-ab28-4866624bed1e",
				"created_at": {
					"seconds": 1533833523,
					"nanos": 537460000
				}
			}]`

			return httpmock.NewStringResponse(200, p), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "workflows", "--project-name", "foo"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET secrets/aaaaaaa")
	}
}
