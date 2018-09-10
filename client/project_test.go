package client

import (
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestInitProjectFromYaml(t *testing.T) {
	yaml := `apiVersion: v1alpha
kind: Project
metadata:
  name: test
spec:
  repository:
    url: "git@github.com:/renderedtext/sem.git"
`
	project, _ := InitProjectFromYaml([]byte(yaml))

	if project.Metadata.Name != "test" {
		t.Errorf("Name is incorrect, got: %s, want: %s.", project.Metadata.Name, "test")
	}

	if project.Spec.Repository.Url != "git@github.com:/renderedtext/sem.git" {
		t.Errorf("Repo Url is incorrect, got: %s, want: %s.", project.Spec.Repository.Url, "git@github.com:/renderedtext/sem.git")
	}
}

func TestProjectToYaml(t *testing.T) {
	project := InitProject("test", "github.com:/renderedtext/sem.git")

	json_body, _ := project.ToJson()

	expected_json_body := `{"apiVersion":"v1alpha","kind":"Project","metadata":{"name":"test"},"spec":{"repository":{"url":"github.com:/renderedtext/sem.git"}}}`

	if string(json_body) != expected_json_body {
		t.Errorf("JSON body is incorrect, got: %s, want: %s.", json_body, expected_json_body)
	}
}

func TestProjectCreate__Success(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	project := InitProject("test", "github.com:/renderedtext/sem.git")
	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/projects",
		func(req *http.Request) (*http.Response, error) {
			jsonbody, _ := project.ToJson()

			return httpmock.NewStringResponse(200, string(jsonbody)), nil
		},
	)

	err := project.Create()

	if err != nil {
		t.Errorf("Expected no errors, got: %s.", err.Error())
	}
}

func TestProjectCreate__ValidationError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/projects",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(422, "Name already taken"), nil
		},
	)

	project := InitProject("test", "github.com:/renderedtext/sem.git")

	err := project.Create()

	expected := `http status 422 with message "Name already taken" received from upstream`

	if err.Error() != expected {
		t.Errorf("Expected an error, got: %s, want: %s.", err, expected)
	}
}
