package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
	assert "github.com/stretchr/testify/assert"
)

func Test__CreateSecret__WithSubcommand__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	content1 := "This is some docker config"
	content2 := "This is some gcloud config"

	ioutil.WriteFile("/tmp/docker", []byte(content1), 0644)
	ioutil.WriteFile("/tmp/gcloud", []byte(content2), 0644)

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1beta/secrets",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{
		"create",
		"secret",
		"-e", "FOO=BAR",
		"-e", "ZEZ=Hello World",
		"--file", "/tmp/docker:.config/docker",
		"--file", "/tmp/gcloud:.config/gcloud",
		"abc",
	})

	RootCmd.Execute()

	file1 := base64.StdEncoding.EncodeToString([]byte(content1))
	file2 := base64.StdEncoding.EncodeToString([]byte(content2))

	expected := fmt.Sprintf(`{"apiVersion":"v1beta","kind":"Secret","metadata":{"name":"abc"},"data":{"env_vars":[{"name":"FOO","value":"BAR"},{"name":"ZEZ","value":"Hello World"}],"files":[{"path":".config/docker","content":"%s"},{"path":".config/gcloud","content":"%s"}]}}`, file1, file2)

	assert.Equal(t, received, expected)
}
