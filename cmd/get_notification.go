package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"
	models "github.com/semaphoreci/cli/api/models"

	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

func NewGetNotificationCmd() *cobra.Command {
	cmd := &cobra.Command{}

	cmd.Use = "notifications [NAME]"
	cmd.Short = "Get notifications."
	cmd.Long = ``
	cmd.Aliases = []string{"notification", "notifs", "notif"}
	cmd.Args = cobra.RangeArgs(0, 1)

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			RunListNotifications(cmd, args)
		} else {
			RunGetNotification(cmd, args)
		}
	}

	return cmd
}

func RunGetNotification(cmd *cobra.Command, args []string) {
	name := args[0]

	c := client.NewNotificationsV1AlphaApi()
	notif, err := c.GetNotification(name)

	utils.Check(err)

	y, err := notif.ToYaml()

	utils.Check(err)

	fmt.Printf("%s", y)
}

func RunListNotifications(cmd *cobra.Command, args []string) {
	c := client.NewNotificationsV1AlphaApi()

	pageSize, _ := cmd.Flags().GetInt32("page-size")
	pageToken, _ := cmd.Flags().GetString("page-token")
	fetchAll := (pageSize == 0) && (pageToken == "")

	var allNotifications []models.NotificationV1Alpha

	for {
		notifList, err := c.ListNotifications(pageSize, pageToken)
		utils.Check(err)

		allNotifications = append(allNotifications, notifList.Notifications...)
		pageToken = notifList.NextPageToken

		if !fetchAll || notifList.NextPageToken == "" {
			break
		}
	}

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	fmt.Fprintln(w, "NAME\tAGE")

	for _, n := range allNotifications {
		updateTime, err := n.Metadata.UpdateTime.Int64()
		utils.Check(err)

		fmt.Fprintf(w, "%s\t%s\n", n.Metadata.Name, utils.RelativeAgeForHumans(updateTime))
	}
	if !fetchAll && pageToken != "" {
		fmt.Fprintf(w, "\nNext page token: %s\n", pageToken)
		fmt.Fprintf(w, "To view next page, run: sem get notifications --page-token %s\n", pageToken)
	}

	if err := w.Flush(); err != nil {
		fmt.Printf("Error flushing when pretty printing notifications: %v\n", err)
	}
}
