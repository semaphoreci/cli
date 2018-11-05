package cmd

import (
	"fmt"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

func NewCreateNotificationCmd() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "notification [NAME]"
	cmd.Short = "Create a notification."
	cmd.Long = ``
	cmd.Aliases = []string{"notification", "notifs", "notif"}
	cmd.Run = RunCreateNotification

	return cmd
}

func RunCreateNotification(cmd *cobra.Command, args []string) {
	name := args[0]

	notif := models.NewNotificationV1Alpha(name)

	notifApi := client.NewNotificationsV1AlphaApi()
	notif, err := notifApi.CreateNotification(notif)

	utils.Check(err)

	fmt.Printf("Notification '%s' created.\n", notif.Metadata.Name)
}
