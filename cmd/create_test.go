package cmd

import (
	"io/ioutil"
	"net/http"
	"testing"

	assert "github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"

	models "github.com/semaphoreci/cli/api/models"
	"github.com/semaphoreci/cli/api/uuid"
)

func Test__CreateProject__FromYaml__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yaml_file := `
apiVersion: v1alpha
kind: Project
metadata:
  name: Test
spec:
  visibility: public
  repository:
    url: "git@github.com:/semaphoreci/cli.git"
    integration_type: github_token
`

	yaml_file_path := "/tmp/project.yaml"

	ioutil.WriteFile(yaml_file_path, []byte(yaml_file), 0644)

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/projects",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{"create", "-f", yaml_file_path})
	RootCmd.Execute()

	expected := `{"apiVersion":"v1alpha","kind":"Project","metadata":{"name":"Test"},"spec":{"visibility":"public","repository":{"url":"git@github.com:/semaphoreci/cli.git","forked_pull_requests":{},"pipeline_file":"","whitelist":{},"integration_type":"github_token"}}}`

	if received != expected {
		t.Errorf("Expected the API to receive POST projects with: %s, got: %s", expected, received)
	}
}

func Test__CreateNotification__FromYaml__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var received *models.NotificationV1Alpha

	endpoint := "https://org.semaphoretext.xyz/api/v1alpha/notifications"

	httpmock.RegisterResponder("POST", endpoint,
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received, _ = models.NewNotificationV1AlphaFromJson(body)

			return httpmock.NewStringResponse(200, "{}"), nil
		},
	)

	RootCmd.SetArgs([]string{"create", "-f", "../fixtures/notification.yml"})
	RootCmd.Execute()

	assert.Equal(t, received.Metadata.Name, "test")

	rule := received.Spec.Rules[0]

	assert.Equal(t, rule.Name, "Rule #1")
	assert.Equal(t, rule.Filter.Projects, []string{"cli"})
	assert.Equal(t, rule.Filter.Branches, []string{"master"})
	assert.Equal(t, rule.Filter.Pipelines, []string{"semaphore.yml"})

	assert.Equal(t, rule.Notify.Slack.Endpoint, "https://hooks.slack.com/asdasdasd/sada/sdas/da")
}

func Test__CreateSecret__FromYaml__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yaml_file := `
apiVersion: v1beta
kind: Secret
metadata:
  name: Test
data:
  env_vars:
  - value: A
    name: B
  files:
  - path: "a.txt"
    content: "21313123"
`

	yaml_file_path := "/tmp/project.yaml"

	ioutil.WriteFile(yaml_file_path, []byte(yaml_file), 0644)

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1beta/secrets",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{"create", "-f", yaml_file_path})
	RootCmd.Execute()

	expected := `{"apiVersion":"v1beta","kind":"Secret","metadata":{"name":"Test"},"data":{"env_vars":[{"name":"B","value":"A"}],"files":[{"path":"a.txt","content":"21313123"}]}}`

	if received != expected {
		t.Errorf("Expected the API to receive POST secret with: %s, got: %s", expected, received)
	}
}

func Test__CreateDashboard__FromYaml__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yaml_file := `
apiVersion: v1alpha
kind: Dashboard
metadata:
  name: Test
  title: "Test Something"
spec:
  widgets:
    - name: "Workflows"
      type: list
      filters:
         github_uid: "{{ github_uid }}"
`

	yaml_file_path := "/tmp/project.yaml"

	ioutil.WriteFile(yaml_file_path, []byte(yaml_file), 0644)

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/dashboards",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{"create", "-f", yaml_file_path})
	RootCmd.Execute()

	expected := `{"apiVersion":"v1alpha","kind":"Dashboard","metadata":{"name":"Test","title":"Test Something"},"spec":{"widgets":[{"name":"Workflows","type":"list","filters":{"github_uid":"{{ github_uid }}"}}]}}`

	if received != expected {
		t.Errorf("Expected the API to receive POST dashbord with: %s, got: %s", expected, received)
	}
}

func Test__CreateDashboard__WithSubcommand__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/dashboards",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{"create", "dash", "abc"})
	RootCmd.Execute()

	expected := `{"apiVersion":"v1alpha","kind":"Dashboard","metadata":{"name":"abc"},"spec":{}}`

	if received != expected {
		t.Errorf("Expected the API to receive POST dashboard with: %s, got: %s", expected, received)
	}
}

func Test__CreateAgentType__FromYaml__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	yaml_file := `
apiVersion: v1alpha
kind: SelfHostedAgentType
metadata:
  name: s1-testing-from-yaml
`

	yaml_file_path := "/tmp/agent_type.yaml"

	ioutil.WriteFile(yaml_file_path, []byte(yaml_file), 0644)

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/self_hosted_agent_types",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{"create", "-f", yaml_file_path})
	RootCmd.Execute()

	expected := `{"apiVersion":"v1alpha","kind":"SelfHostedAgentType","metadata":{"name":"s1-testing-from-yaml"},"status":{}}`

	if received != expected {
		t.Errorf("Expected the API to receive POST self_hosted_agent_types with: %s, got: %s", expected, received)
	}
}

func Test__CreateAgentType__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/self_hosted_agent_types",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{"create", "agent_type", "s1-testing"})
	RootCmd.Execute()

	expected := `{"apiVersion":"v1alpha","kind":"SelfHostedAgentType","metadata":{"name":"s1-testing"},"status":{}}`

	if received != expected {
		t.Errorf("Expected the API to receive POST self_hosted_agent_types with: %s, got: %s", expected, received)
	}
}

func Test__CreateDeploymentTarget__FromYaml__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	uuid.Mock()
	defer uuid.Unmock()

	yaml_file := `
apiVersion: v1alpha
kind: DeploymentTarget
metadata:
  name: dt-name-from-yaml
  organization_id: org-id
  project_id: prj-id
  url: www.semaphore.xyz
  description: dt-description
spec:
  bookmark_parameter1: book1
`

	yaml_file_path := "/tmp/create_dt.yaml"

	ioutil.WriteFile(yaml_file_path, []byte(yaml_file), 0644)

	received := ""

	httpmock.RegisterResponder(http.MethodPost, "https://org.semaphoretext.xyz/api/v1alpha/deployment_targets",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{"create", "-f", yaml_file_path})
	RootCmd.Execute()

	expected := `{"id":"","name":"dt-name-from-yaml","project_id":"prj-id","organization_id":"org-id","description":"dt-description","url":"www.semaphore.xyz","state":"","state_message":"","subject_rules":null,"object_rules":null,"active":false,"bookmark_parameter1":"book1","bookmark_parameter2":"","bookmark_parameter3":"","unique_token":"00020406-090b-4e10-9315-181a1c1e2022"}`
	if received != expected {
		t.Errorf("Expected the API to receive POST deployment_targets: %s, got: %s", expected, received)
	}
}
