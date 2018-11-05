package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"

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

	notifList, err := c.ListNotifications()

	utils.Check(err)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

	fmt.Fprintln(w, "NAME")

	for _, n := range notifList.Notifications {
		updateTime, err := n.Metadata.UpdateTime.Int64()

		utils.Check(err)

		fmt.Fprintf(w, "%s\t%s\n", n.Metadata.Name, utils.RelativeAgeForHumans(updateTime))
	}

	w.Flush()
}
