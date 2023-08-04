package cmd

import (
	"bytes"
	"net/http"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
	assert "github.com/stretchr/testify/assert"
)

func Test__Connect__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var outputBuffer bytes.Buffer

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/projects",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(200, "[]"), nil
		},
	)

	RootCmd.SetArgs([]string{"connect", "org.semaphoretext.xyz", "abc"})
	RootCmd.SetOutput(&outputBuffer)
	RootCmd.Execute()

	assert.Equal(t, outputBuffer.String(), "connected to org.semaphoretext.xyz\n")
}

func Test__Connect__Response401(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var outputBuffer bytes.Buffer

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/projects",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(401, "Unauthorized"), nil
		},
	)

	RootCmd.SetArgs([]string{"connect", "org.semaphoretext.xyz", "abc"})
	RootCmd.SetOutput(&outputBuffer)

	// handle exit
	assert.Panics(t, func() {
		RootCmd.Execute()
	}, "exit 1")

	assert.Equal(
		t,
		outputBuffer.String(),
		"http status 401 with message \"Unauthorized\" received from upstream")
}
