package cmd

import (
	"fmt"
	"strings"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"
	utils "github.com/semaphoreci/cli/cmd/utils"

	"github.com/spf13/cobra"
)

func NewCreateNotificationCmd() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "notification [NAME]"
	cmd.Short = "Create a notification."
	cmd.Long = ``
	cmd.Aliases = []string{"notifications", "notifs", "notif"}
	cmd.Args = cobra.ExactArgs(1)
	cmd.Run = RunCreateNotification

	cmd.Flags().String("projects", "", "Filter for project names")
	cmd.Flags().String("pipelines", "", "Filter for pipeline file")
	cmd.Flags().String("branches", "", "Filter for branch names")
	cmd.Flags().String("tags", "", "Filter for tag names")

	cmd.Flags().String("slack-channels", "", "Slack channels where to send notifications")
	cmd.Flags().String("slack-endpoint", "", "Slack webhook endpoint")

	cmd.Flags().String("webhook-endpoint", "", "Webhook endpoint")
	cmd.Flags().String("webhook-secret", "", "Webhook secret")

	return cmd
}

const errNotificationWithoutProject = `Specify at least one project that sends notifications.

Example:

  sem create notification my-notif --projects "cli,webapp"
`

const errNotificationWithoutSlackChannels = `Specify at least one slack channel where to send notifications.

Example:

  sem create notification my-notif --slack-channels "#general,#devops"
`

const errNotificationWithoutEndpoint = `Specify the slack and/or webhook endpoint where to send notificaitons.

Example:

  sem create notification my-notif \
    --slack-endpoint "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX" \
    --webhook-endpoint "https://example.com/postrequest"
`

func RunCreateNotification(cmd *cobra.Command, args []string) {
	name := args[0]

	projects, err := utils.CSVFlag(cmd, "projects")
	utils.Check(err)

	branches, err := utils.CSVFlag(cmd, "branches")
	utils.Check(err)

	tags, err := utils.CSVFlag(cmd, "tags")
	utils.Check(err)

	pipelines, err := utils.CSVFlag(cmd, "pipelines")
	utils.Check(err)

	slackChannels, err := utils.CSVFlag(cmd, "slack-channels")
	utils.Check(err)

	slackEndpoint, err := cmd.Flags().GetString("slack-endpoint")
	utils.Check(err)

	webhookEndpoint, err := cmd.Flags().GetString("webhook-endpoint")
	utils.Check(err)

	webhookSecret, err := cmd.Flags().GetString("webhook-secret")
	utils.Check(err)

	if len(projects) == 0 {
		utils.Fail(errNotificationWithoutProject)
	}

	if slackEndpoint == "" && webhookEndpoint == "" {
		utils.Fail(errNotificationWithoutEndpoint)
	}

	filter := models.NotificationV1AlphaSpecRuleFilter{}
	filter.Projects = projects
	filter.Branches = branches
	filter.Tags = tags
	filter.Pipelines = pipelines

	notify := models.NotificationV1AlphaSpecRuleNotify{}
	notify.Slack.Channels = slackChannels
	notify.Slack.Endpoint = slackEndpoint
	notify.Webhook.Endpoint = webhookEndpoint
	notify.Webhook.Secret = webhookSecret

	ruleName := fmt.Sprintf(
		"Send notifications for %s", strings.Join(projects, ", "))

	rule := models.NotificationV1AlphaSpecRule{
		Name:   ruleName,
		Filter: filter,
		Notify: notify,
	}

	notif := models.NewNotificationV1Alpha(name)
	notif.Spec.Rules = append(notif.Spec.Rules, rule)

	notifApi := client.NewNotificationsV1AlphaApi()
	notif, err = notifApi.CreateNotification(notif)

	utils.Check(err)

	fmt.Printf("Notification '%s' created.\n", notif.Metadata.Name)
}
