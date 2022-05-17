package cmd

import (
	"io/ioutil"
	"net/http"
	"testing"

	models "github.com/semaphoreci/cli/api/models"
	httpmock "gopkg.in/jarcoal/httpmock.v1"

	assert "github.com/stretchr/testify/assert"
)

func Test__CreateNotification(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var received *models.NotificationV1Alpha

	endpoint := "https://org.semaphoretext.xyz/api/v1alpha/notifications"

	httpmock.RegisterResponder("POST", endpoint,
		func(req *http.Request) (*http.Response, error) {
			body, _ := ioutil.ReadAll(req.Body)

			received, _ = models.NewNotificationV1AlphaFromJson(body)

			return httpmock.NewStringResponse(200, "{}"), nil
		},
	)

	RootCmd.SetArgs([]string{
		"create",
		"notification",
		"aaa",
		"--projects", "cli, test1",
		"--branches", "master,staging",
		"--pipelines", ".semaphore/semaphore.yml",
		"--slack-channels", "#product",
		"--slack-endpoint", "https://dasdasdasd/sa/das/da/sdas",
		"--webhook-endpoint", "https://dasdasdasd/sa",
		"--webhook-secret", "aaa-webhook-secret",
	})

	RootCmd.Execute()

	assert.Equal(t, received.Metadata.Name, "aaa")

	assert.Equal(t, received.Spec.Rules[0].Filter.Projects, []string{
		"cli",
		"test1",
	})

	assert.Equal(t, received.Spec.Rules[0].Filter.Branches, []string{
		"master",
		"staging",
	})

	assert.Equal(t, received.Spec.Rules[0].Filter.Pipelines, []string{
		".semaphore/semaphore.yml",
	})

	assert.Equal(t, received.Spec.Rules[0].Notify.Slack.Channels, []string{
		"#product",
	})

	assert.Equal(t, received.Spec.Rules[0].Notify.Slack.Endpoint,
		"https://dasdasdasd/sa/das/da/sdas")

	assert.Equal(t, received.Spec.Rules[0].Notify.Webhook.Endpoint,
		"https://dasdasdasd/sa")

	assert.Equal(t, received.Spec.Rules[0].Notify.Webhook.Secret,
		"aaa-webhook-secret")
}
