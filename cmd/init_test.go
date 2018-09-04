package cmd

import (
	"io/ioutil"
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestRunInit(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/projects",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received := body

			return httpmock.NewStringResponse(200, string(received)), nil
		},
	)

	RunInit(InitCmd, []string{})

	expected := `{"metadata":{"name":"something"},"spec":{"repository":{"url":"git@github.com:/renderedtext/something.git"}}}`

	if received == expected {
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
	_, err := ConstructProjectName("github.com/renderedtext/test")

	if err == nil {
		t.Errorf("Expected error for unrecognized format.")
	}
}
