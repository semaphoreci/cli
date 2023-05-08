package cmd

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func Test__RebuildPipeline__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/pipelines/494b76aa-f3f0-4ecf-b5ef-c389591a01be/partial_rebuild",
		func(req *http.Request) (*http.Response, error) {

			received = "true"

			return httpmock.NewStringResponse(200, "message"), nil
		},
	)

	RootCmd.SetArgs([]string{"rebuild", "pipeline", "494b76aa-f3f0-4ecf-b5ef-c389591a01be"})
	RootCmd.Execute()

	expected := "true"

	if received != expected {
		t.Errorf("Expected the API to receive POST pipelines with: %s, got: %s", expected, received)
	}
}

func Test__RebuildDeploymentTarget__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	targetId := "494b76aa-f3f0-4ecf-b5ef-c389591a01be"
	patchURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets/%s/on", targetId)
	httpmock.RegisterResponder(http.MethodPatch, patchURL,
		func(req *http.Request) (*http.Response, error) {
			received = true

			target := `{
				"target_id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
			  "cordoned": false
	  		}	
	  		`

			return httpmock.NewStringResponse(200, target), nil
		},
	)

	RootCmd.SetArgs([]string{"rebuild", "target", "494b76aa-f3f0-4ecf-b5ef-c389591a01be"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive PATCH deployment_targets/:id")
}
