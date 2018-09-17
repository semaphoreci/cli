package cmd

import (
	"io/ioutil"
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
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
  repository:
    url: "git@github.com:/semaphoreci/cli.git"
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

	expected := `{"apiVersion":"v1alpha","kind":"Project","metadata":{"name":"Test"},"spec":{"repository":{"url":"git@github.com:/semaphoreci/cli.git"}}}`

	if received != expected {
		t.Errorf("Expected the API to receive POST projects with: %s, got: %s", expected, received)
	}
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

func Test__CreateSecret__WithSubcommand__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1beta/secrets",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{"create", "secret", "abc"})
	RootCmd.Execute()

	expected := `{"apiVersion":"v1beta","kind":"Secret","metadata":{"name":"abc"},"data":{"env_vars":null,"files":null}}`

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
