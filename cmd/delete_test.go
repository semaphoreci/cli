package cmd

import (
	"fmt"
	"net/http"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
	"github.com/semaphoreci/cli/api/uuid"
	assert "github.com/stretchr/testify/assert"
)

func TestDeleteProject__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("DELETE", "https://org.semaphoretext.xyz/api/v1alpha/projects/test-prj",
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	RootCmd.SetArgs([]string{"delete", "project", "test-prj"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive DELETE test-prj")
	}
}

func TestDeleteSecret__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("DELETE", "https://org.semaphoretext.xyz/api/v1beta/secrets/test-secret",
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	RootCmd.SetArgs([]string{"delete", "secret", "test-secret"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive DELETE test-prj")
	}
}

func TestDeleteDashboardCmd__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("DELETE", "https://org.semaphoretext.xyz/api/v1alpha/dashboards/test-dash",
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	RootCmd.SetArgs([]string{"delete", "dash", "test-dash"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive DELETE test dash")
	}
}

func TestDeleteAgentTypeCmd__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("DELETE", "https://org.semaphoretext.xyz/api/v1alpha/self_hosted_agent_types/s1-testing",
		func(req *http.Request) (*http.Response, error) {
			received = true

			return httpmock.NewStringResponse(200, ""), nil
		},
	)

	RootCmd.SetArgs([]string{"delete", "agent_type", "s1-testing"})
	RootCmd.Execute()

	if received == false {
		t.Error("Expected the API to receive DELETE agent_type s1-testing")
	}
}

func Test__DeleteDeploymentTarget__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	uuid.Mock()
	defer uuid.Unmock()

	received := false

	targetId := "494b76aa-f3f0-4ecf-b5ef-c389591a01be"
	unique_token, _ := uuid.NewUUID()
	deleteURL := fmt.Sprintf("https://org.semaphoretext.xyz/api/v1alpha/deployment_targets/%s?unique_token=%s", targetId, unique_token)
	httpmock.RegisterResponder(http.MethodDelete, deleteURL,
		func(req *http.Request) (*http.Response, error) {
			received = true

			p := `{
					"id": "494b76aa-f3f0-4ecf-b5ef-c389591a01be",
					"name": "dep target test",
					"url": "https://semaphoreci.xyz/target",
					"project_id": "proj_id"
			}
			`

			return httpmock.NewStringResponse(200, p), nil
		},
	)

	RootCmd.SetArgs([]string{"delete", "dt", targetId})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive DELETE deployment_targets/:id")
}
