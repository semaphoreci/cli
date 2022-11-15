package cmd

import (
	"net/http"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"
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
