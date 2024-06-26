package cmd

import (
	"net/http"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
)

const Events = `{
  "events": [
		{"event": "job_started", "timestamp": 1624541916},
		{"event": "cmd_started", "timestamp": 1624541916, "directive": "Exporting environment variables"},
		{"event": "cmd_output", "timestamp": 1624541916, "output": "Exporting VAR1"},
		{"event": "cmd_output", "timestamp": 1624541916, "output": "Exporting VAR2"},
		{"event": "cmd_finished", "timestamp": 1624541916, "directive": "Exporting environment variables", "exit_code": 0},
		{"event": "job_finished", "timestamp": 1624541916, "result": "passed"}
	]
}
`

func TestLogs__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/logs/job-123",
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, Events), nil
		},
	)

	RootCmd.SetArgs([]string{"logs", "job-123", "-v"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive GET /logs/job-123")
	}
}
