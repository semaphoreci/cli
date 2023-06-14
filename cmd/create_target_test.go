package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/semaphoreci/cli/api/uuid"
	assert "github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func Test__CreateDeploymentTarget__WithSubcommand__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	uuid.Mock()
	defer uuid.Unmock()

	content1 := "This is some docker config"
	content2 := "This is some gcloud config"

	ioutil.WriteFile("/tmp/docker", []byte(content1), 0644)
	ioutil.WriteFile("/tmp/gcloud", []byte(content2), 0644)

	received := ""

	httpmock.RegisterResponder("POST", "https://org.semaphoretext.xyz/api/v1alpha/deployment_targets",
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received = string(body)
			fmt.Printf("received:'%s'\n", received)

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{
		"create",
		"dt",
		"-i", "00000000-0000-0000-000000000000",
		"-e", "FOO=BAR",
		"-e", "ZEZ=Hello World",
		"--file", "/tmp/docker:.config/docker",
		"--file", "/tmp/gcloud:.config/gcloud",
		"abc",
	})

	RootCmd.Execute()

	file1 := base64.StdEncoding.EncodeToString([]byte(content1))
	file2 := base64.StdEncoding.EncodeToString([]byte(content2))

	expected := fmt.Sprintf(`{"id":"","name":"abc","project_id":"00000000-0000-0000-000000000000","organization_id":"","description":"","url":"","state":"","state_message":"","subject_rules":null,"object_rules":null,"active":false,"bookmark_parameter1":"","bookmark_parameter2":"","bookmark_parameter3":"","env_vars":[{"name":"FOO","value":"BAR"},{"name":"ZEZ","value":"Hello World"}],"files":[{"path":".config/docker","content":"%s"},{"path":".config/gcloud","content":"%s"}],"unique_token":"00020406-090b-4e10-9315-181a1c1e2022"}`, file1, file2)

	assert.Equal(t, received, expected)
}
