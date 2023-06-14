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

			return httpmock.NewStringResponse(200, received), nil
		},
	)

	RootCmd.SetArgs([]string{
		"create",
		"dt",
		"-i", "00000000-0000-0000-000000000000",
		"-e", "FOO=BAR",
		"--env", "ZEZ=Hello World",
		"--file", "/tmp/docker:.config/docker",
		"-f", "/tmp/gcloud:.config/gcloud",
		"-s", "any",
		"--subject-rule", "user,mock_user_321",
		"-s", "role,contributor",
		"--object-rule", "branch,exact,main",
		"-o", `tag,regex,.*feat.*`,
		"-b", "book 1",
		"--url", "mock_url_321.zyx",
		"abc",
	})

	RootCmd.Execute()

	file1 := base64.StdEncoding.EncodeToString([]byte(content1))
	file2 := base64.StdEncoding.EncodeToString([]byte(content2))

	expected := fmt.Sprintf(`{"id":"","name":"abc","project_id":"00000000-0000-0000-000000000000","organization_id":"","description":"","url":"mock_url_321.zyx","state":"","state_message":"","subject_rules":[{"type":"ANY"},{"type":"USER","git_login":"mock_user_321"},{"type":"ROLE","subject_id":"contributor"}],"object_rules":[{"type":"branch","match_mode":"exact","pattern":"main"},{"type":"tag","match_mode":"regex","pattern":".*feat.*"}],"active":false,"bookmark_parameter1":"book 1","bookmark_parameter2":"","bookmark_parameter3":"","env_vars":[{"name":"FOO","value":"BAR"},{"name":"ZEZ","value":"Hello World"}],"files":[{"path":".config/docker","content":"%s"},{"path":".config/gcloud","content":"%s"}],"unique_token":"00020406-090b-4e10-9315-181a1c1e2022"}`, file1, file2)

	assert.Equal(t, received, expected)
}
