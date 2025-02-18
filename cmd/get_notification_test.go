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

func Test__ListNotifications__WithAllFlag(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	page1Called := false
	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/notifications",
		func(req *http.Request) (*http.Response, error) {
			page1Called = true

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

			response := fmt.Sprintf(`{
				"notifications": [%s,%s],
				"next_page_token": "page2"
			}`, n1, n2)

			return httpmock.NewStringResponse(200, response), nil
		},
	)

	page2Called := false
	httpmock.RegisterResponder("GET", "https://org.semaphoretext.xyz/api/v1alpha/notifications?page_token=page2",
		func(req *http.Request) (*http.Response, error) {
			page2Called = true

			n3 := `{
				"metadata": {
					"name": "notif3",
					"update_time": "126"
				}
			}`

			response := fmt.Sprintf(`{
				"notifications": [%s],
				"next_page_token": ""
			}`, n3)

			return httpmock.NewStringResponse(200, response), nil
		},
	)

	RootCmd.SetArgs([]string{"get", "notifications"})
	RootCmd.Execute()

	assert.True(t, page1Called, "Expected the API to receive GET notifications for first page")
	assert.True(t, page2Called, "Expected the API to receive GET notifications for second page")
}
