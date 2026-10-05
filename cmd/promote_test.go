package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
	"github.com/spf13/pflag"
	assert "github.com/stretchr/testify/assert"
)

// Cobra keeps flag values between executions, so each test starts from the
// defaults.
func resetPromoteFlags() {
	promoteCmd.Flags().Set("override", "false")

	if params, ok := promoteCmd.Flags().Lookup("param").Value.(pflag.SliceValue); ok {
		params.Replace([]string{})
	}
}

func Test__Promote__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var received map[string]interface{}

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/promotions",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &received)

			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	resetPromoteFlags()
	RootCmd.SetArgs([]string{"promote", "494b76aa-f3f0-4ecf-b5ef-c389591a01be", "Deploy to production"})
	RootCmd.Execute()

	assert.Equal(t, map[string]interface{}{
		"pipeline_id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
		"name":        "Deploy to production",
		"override":    false,
	}, received)
}

func Test__Promote__WithParametersAndOverride(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var received map[string]interface{}

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/promotions",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &received)

			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	resetPromoteFlags()
	RootCmd.SetArgs([]string{
		"promote", "494b76aa-f3f0-4ecf-b5ef-c389591a01be", "Deploy",
		"--param", "REGION=eu-west-1",
		"-p", "MESSAGE=key=value with spaces",
		"--override",
	})
	RootCmd.Execute()

	assert.Equal(t, map[string]interface{}{
		"pipeline_id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
		"name":        "Deploy",
		"override":    true,
		"REGION":      "eu-west-1",
		"MESSAGE":     "key=value with spaces",
	}, received)
}

func Test__Promote__InvalidParameter(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	called := false

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/promotions",
		func(req *http.Request) (*http.Response, error) {
			called = true

			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	resetPromoteFlags()
	RootCmd.SetArgs([]string{"promote", "494b76aa-f3f0-4ecf-b5ef-c389591a01be", "Deploy", "--param", "REGION"})

	assert.PanicsWithValue(t, "exit 1", func() { RootCmd.Execute() })
	assert.False(t, called, "the API must not be called when a --param is malformed")
}

func Test__Promote__ReservedParameterName(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	called := false

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/promotions",
		func(req *http.Request) (*http.Response, error) {
			called = true

			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	resetPromoteFlags()
	RootCmd.SetArgs([]string{"promote", "494b76aa-f3f0-4ecf-b5ef-c389591a01be", "Deploy", "--param", "name=other"})

	assert.PanicsWithValue(t, "exit 1", func() { RootCmd.Execute() })
	assert.False(t, called, "the API must not be called when a parameter uses a reserved name")
}

func Test__Promote__Response4xx(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/promotions",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(400, "Triggering promotion with deployment target failed: object not allowed"), nil
		},
	)

	resetPromoteFlags()
	RootCmd.SetArgs([]string{"promote", "494b76aa-f3f0-4ecf-b5ef-c389591a01be", "Deploy"})

	assert.PanicsWithValue(t, "exit 1", func() { RootCmd.Execute() })
}
