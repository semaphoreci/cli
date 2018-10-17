package cmd

import (
	"io/ioutil"
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func Test__StopPipeline__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := ""

	httpmock.RegisterResponder("PATCH", "https://org.semaphoretext.xyz/api/v1alpha/pipelines/494b76aa-f3f0-4ecf-b5ef-c389591a01be",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, "message"), nil
		},
	)

	RootCmd.SetArgs([]string{"stop", "pipeline", "494b76aa-f3f0-4ecf-b5ef-c389591a01be"})
	RootCmd.Execute()

	expected := "{\"terminate_request\": true}"

	if received != expected {
		t.Errorf("Expected the API to receive PATCH pipelines with: %s, got: %s", expected, received)
	}
}
