package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestRunInit__NoParams(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/projects",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, string(received)), nil
		},
	)

	cmd := InitCmd()
	cmd.Execute()

	expected := `{"apiVersion":"v1alpha","kind":"Project","metadata":{"name":"something"},"spec":{"repository":{"url":"git@github.com:/renderedtext/something.git","run_on":["branches","tags"],"forked_pull_requests":{},"pipeline_file":"","whitelist":{}}}}`

	if received != expected {
		t.Errorf("Expected the API to receive project create req with '%s' instead '%s'.", expected, received)
	}
}

func TestRunInit__NameParamPassed(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/projects",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, string(received)), nil
		},
	)

	cmd := InitCmd()

	cmd.SetArgs([]string{
		`--project-name=another-name`,
	})

	cmd.Execute()

	expected := `{"apiVersion":"v1alpha","kind":"Project","metadata":{"name":"another-name"},"spec":{"repository":{"url":"git@github.com:/renderedtext/something.git","run_on":["branches","tags"],"forked_pull_requests":{},"pipeline_file":"","whitelist":{}}}}`

	fmt.Print(expected)
	fmt.Print(received)

	if received != expected {
		t.Errorf("Expected the API to receive project create req, want: %s got: %s.", expected, received)
	}
}

func TestRunInit__RepoUrlParamPassed(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/projects",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, string(received)), nil
		},
	)

	cmd := InitCmd()

	cmd.SetArgs([]string{
		`--repo-url=git@github.com:/renderedtext/a.git`,
	})

	cmd.Execute()

	expected := `{"apiVersion":"v1alpha","kind":"Project","metadata":{"name":"a"},"spec":{"repository":{"url":"git@github.com:/renderedtext/a.git","run_on":["branches","tags"],"forked_pull_requests":{},"pipeline_file":"","whitelist":{}}}}`

	fmt.Print(expected)
	fmt.Print(received)

	if received != expected {
		t.Errorf("Expected the API to receive project create req, want: %s got: %s.", expected, received)
	}
}

func TestConstructProjectName__GitFormat(t *testing.T) {
	name, _ := ConstructProjectName("git@github.com:/renderedtext/test.git")

	if name != "test" {
		t.Errorf("Name is incorrect, got: %s, want: %s.", name, "test")
	}

	name, _ = ConstructProjectName("git@github.com:/renderedtext/test")

	if name != "test" {
		t.Errorf("Name is incorrect, got: %s, want: %s.", name, "test")
	}
}

func TestConstructProjectName__HttpFormat(t *testing.T) {
	name, _ := ConstructProjectName("https://semaphore@bitbucket.org/semaphoreci/test.git")

	if name != "test" {
		t.Errorf("Name is incorrect, got: %s, want: %s.", name, "test")
	}

	name, _ = ConstructProjectName("https://semaphore@bitbucket.org/semaphoreci/test")

	if name != "test" {
		t.Errorf("Name is incorrect, got: %s, want: %s.", name, "test")
	}
}

func TestConstructProjectName__InvalidFormat(t *testing.T) {
	_, err := ConstructProjectName("github.com/renderedtext")

	if err == nil {
		t.Errorf("Expected error for unrecognized format.")
	}
}
