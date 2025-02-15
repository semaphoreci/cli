package cmd

import (
	"fmt"
	"net/http"
	"testing"

	httpmock "github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func Test__ListNotifications__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/notifications",
		func(req *http.Request) (*http.Response, error) {
			received = true

			n1 := `{
				"metadata": {
					"name": "notif1",
					"update_time": "124"
				}
			}`

			n2 := `{
				"metadata": {
					"name": "notif2",
					"update_time": "125"
				}
			}`

			notifications := fmt.Sprintf(`{"notifications":[%s,%s]}`, n1, n2)

			return httpmock.NewStringResponse(200, notifications), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "notifications"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET notifications")
}

func Test__ListNotifications__WithPageToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/notifications?page_size=50&page_token=abc123",
		func(req *http.Request) (*http.Response, error) {
			received = true

			response := `{
				"notifications": [{
					"metadata": {
						"name": "notif2",
						"update_time": "125"
					}
				}]
			}`

			return httpmock.NewStringResponse(200, response), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "notifications", "--page-token", "abc123", "--page-size", "50"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET notifications with page token")
}

func Test__GetNotification__Response200(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	received := false

	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/notifications/test-notif",
		func(req *http.Request) (*http.Response, error) {
			received = true

			notification := `{
				"metadata": {
					"name": "test-notif",
					"update_time": "124"
				}
			}`

			return httpmock.NewStringResponse(200, notification), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "notifications", "test-notif"})
	RootCmd.Execute()

	assert.True(t, received, "Expected the API to receive GET notifications/test-notif")
}
