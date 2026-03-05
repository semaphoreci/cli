package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

// resetRunTaskFlags clears flag state between tests because cobra's global
// command tree persists flag values across RootCmd.Execute() calls.
func resetRunTaskFlags() {
	for _, name := range []string{"branch", "tag", "pipeline-file"} {
		f := runTaskCmd.Flags().Lookup(name)
		f.Value.Set("")
		f.Changed = false
	}

	if f := runTaskCmd.Flags().Lookup("param"); f != nil {
		if sv, ok := f.Value.(pflag.SliceValue); ok {
			sv.Replace([]string{})
		}
		f.Changed = false
	}
}

func Test__RunTask__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c/run_now",
		func(req *http.Request) (*http.Response, error) {
			received = true

			resp := `{"workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e"}`
			return httpmock.NewStringResponse(200, resp), nil
		},
	)

	resetRunTaskFlags()
	RootCmd.SetArgs([]string{"run", "task", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive POST /tasks/:id/run_now")
}

func Test__RunTask__NoFlags(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var receivedBody string

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c/run_now",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			receivedBody = string(body)

			resp := `{"workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e"}`
			return httpmock.NewStringResponse(200, resp), nil
		},
	)

	resetRunTaskFlags()
	RootCmd.SetArgs([]string{"run", "task", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c"})
	RootCmd.Execute()

	assert.Equal(t, "{}", receivedBody)
}

func Test__RunTask__WithBranch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var receivedBody map[string]interface{}

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c/run_now",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &receivedBody)

			resp := `{"workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e"}`
			return httpmock.NewStringResponse(200, resp), nil
		},
	)

	resetRunTaskFlags()
	RootCmd.SetArgs([]string{"run", "task", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c", "--branch", "main"})
	RootCmd.Execute()

	assert.NotNil(t, receivedBody)
	ref := receivedBody["reference"].(map[string]interface{})
	assert.Equal(t, "BRANCH", ref["type"])
	assert.Equal(t, "main", ref["name"])
}

func Test__RunTask__WithTag(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var receivedBody map[string]interface{}

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c/run_now",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &receivedBody)

			resp := `{"workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e"}`
			return httpmock.NewStringResponse(200, resp), nil
		},
	)

	resetRunTaskFlags()
	RootCmd.SetArgs([]string{"run", "task", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c", "--tag", "v1.0"})
	RootCmd.Execute()

	assert.NotNil(t, receivedBody)
	ref := receivedBody["reference"].(map[string]interface{})
	assert.Equal(t, "TAG", ref["type"])
	assert.Equal(t, "v1.0", ref["name"])
}

func Test__RunTask__WithPipelineFile(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var receivedBody map[string]interface{}

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c/run_now",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &receivedBody)

			resp := `{"workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e"}`
			return httpmock.NewStringResponse(200, resp), nil
		},
	)

	resetRunTaskFlags()
	RootCmd.SetArgs([]string{"run", "task", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c", "--pipeline-file", ".semaphore/custom.yml"})
	RootCmd.Execute()

	assert.NotNil(t, receivedBody)
	assert.Equal(t, ".semaphore/custom.yml", receivedBody["pipeline_file"])
}

func Test__RunTask__WithParams(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var receivedBody map[string]interface{}

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c/run_now",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &receivedBody)

			resp := `{"workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e"}`
			return httpmock.NewStringResponse(200, resp), nil
		},
	)

	resetRunTaskFlags()
	RootCmd.SetArgs([]string{"run", "task", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c", "--param", "ENV=staging", "--param", "REGION=us-east-1"})
	RootCmd.Execute()

	assert.NotNil(t, receivedBody)
	params := receivedBody["parameters"].(map[string]interface{})
	assert.Equal(t, "staging", params["ENV"])
	assert.Equal(t, "us-east-1", params["REGION"])
}

func Test__RunTask__WithAllFlags(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var receivedBody map[string]interface{}

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c/run_now",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &receivedBody)

			resp := `{"workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e"}`
			return httpmock.NewStringResponse(200, resp), nil
		},
	)

	resetRunTaskFlags()
	RootCmd.SetArgs([]string{
		"run", "task", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c",
		"--branch", "develop",
		"--pipeline-file", ".semaphore/custom.yml",
		"--param", "ENV=staging",
	})
	RootCmd.Execute()

	assert.NotNil(t, receivedBody)

	ref := receivedBody["reference"].(map[string]interface{})
	assert.Equal(t, "BRANCH", ref["type"])
	assert.Equal(t, "develop", ref["name"])
	assert.Equal(t, ".semaphore/custom.yml", receivedBody["pipeline_file"])

	params := receivedBody["parameters"].(map[string]interface{})
	assert.Equal(t, "staging", params["ENV"])
}

func Test__RunTask__TasksAlias(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c/run_now",
		func(req *http.Request) (*http.Response, error) {
			received = true

			resp := `{"workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e"}`
			return httpmock.NewStringResponse(200, resp), nil
		},
	)

	resetRunTaskFlags()
	RootCmd.SetArgs([]string{"run", "tasks", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the 'tasks' alias to work for run task")
}

func Test__RunTask__BranchAndTagConflict(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	resetRunTaskFlags()
	RootCmd.SetArgs([]string{"run", "task", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c", "--branch", "main", "--tag", "v1.0"})
	err := RootCmd.Execute()

	assert.NotNil(t, err, "Expected an error when both --branch and --tag are provided")
	assert.Contains(t, err.Error(), "if any flags in the group [branch tag] are set none of the others can be")
}

func Test__RunTask__WithParamContainingEquals(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var receivedBody map[string]interface{}

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/tasks/bb2ba294-d4b3-48bc-90a7-12dd56e9424c/run_now",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(body, &receivedBody)

			resp := `{"workflow_id": "dd4ba294-d4b3-48bc-90a7-12dd56e9424e"}`
			return httpmock.NewStringResponse(200, resp), nil
		},
	)

	resetRunTaskFlags()
	RootCmd.SetArgs([]string{"run", "task", "bb2ba294-d4b3-48bc-90a7-12dd56e9424c", "--param", "CONN=host=db;port=5432"})
	RootCmd.Execute()

	assert.NotNil(t, receivedBody)
	params := receivedBody["parameters"].(map[string]interface{})
	assert.Equal(t, "host=db;port=5432", params["CONN"])
}
