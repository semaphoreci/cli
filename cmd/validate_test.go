package cmd

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
	assert "github.com/stretchr/testify/assert"
)

func TestValidate__FileExistsPathArgument__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var outputBuffer bytes.Buffer

	os.Chdir(t.TempDir())
	os.Mkdir(".semaphore", 0755)
	f, err := os.Create(".semaphore/semaphore.yml")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString("---\n")
	f.Close()

	requestBodySent := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/yaml",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			requestBodySent = string(body)

			return httpmock.NewStringResponse(200, "{\"pipeline_id\": \"\",\"message\":\"YAML definition is valid.\"}"), nil
		},
	)

	RootCmd.SetArgs([]string{"validate", ".semaphore/semaphore.yml"})
	RootCmd.SetOutput(&outputBuffer)
	RootCmd.Execute()

	requestBodyExpected := "{\"yaml_definition\":\"---\\n\"}"

	if requestBodyExpected != requestBodySent {
		t.Errorf("Expected the API to receive YAML with: %s, got: %s", requestBodyExpected, requestBodySent)
	}

	assert.Equal(t, outputBuffer.String(), "")
}

func TestValidate__FileMissingPathArgument(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	os.Chdir(t.TempDir())

	RootCmd.SetArgs([]string{"validate", "semaphore.yml"})
	assert.Panics(t, func() {
		RootCmd.Execute()
	}, "exit 1")

	// Can't check output as utils.CheckWithMessage() writes directly to stderr
}
