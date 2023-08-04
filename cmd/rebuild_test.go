package cmd

import (
	"net/http"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
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
