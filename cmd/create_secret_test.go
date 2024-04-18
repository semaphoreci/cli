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

func Test__CreateSecret__WithProjectID__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := ""

	dash := `{
		"metadata":{
			"name":"hello",
			"id":"bb2ba294-d4b3-48bc-90a7-12dd56e9424b",
		}
	}`

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/projects/hello",
		func(req *http.Request) (*http.Response, error) {
			fmt.Println("GET /api/v1alpha/projects/hello")
			return httpmock.NewStringResponse(200, dash), nil
		},
	)

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1/projects/hello/secrets",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	content1 := "This is some docker config"
	content2 := "This is some gcloud config"

	ioutil.WriteFile("/tmp/docker", []byte(content1), 0644)
	ioutil.WriteFile("/tmp/gcloud", []byte(content2), 0644)

	// flags for env vars and projects stay the same as in the previous test

	RootCmd.SetArgs([]string{
		"create",
		"secret",
		"-i", "hello",
		"projectABC",
	})


	RootCmd.Execute()

	
	file1 := base64.StdEncoding.EncodeToString([]byte(content1))
	file2 := base64.StdEncoding.EncodeToString([]byte(content2))	// We do not expect project_id_or_name to be received in the body of the request, as it's set by the grpc-gateway from the URL
	expected := fmt.Sprintf(`{"apiVersion":"v1","kind":"ProjectSecret","metadata":{"name":"projectABC"},"data":{"env_vars":[{"name":"FOO","value":"BAR"},{"name":"ZEZ","value":"Hello World"}],"files":[{"path":".config/docker","content":"%s"},{"path":".config/gcloud","content":"%s"}]}}`, file1, file2)

	assert.Equal(t, expected, received)
}
