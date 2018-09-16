package cmd_get

import (
	"fmt"
	"os"
	"text/tabwriter"

	client "github.com/semaphoreci/cli/api/client"

	"github.com/semaphoreci/cli/cmd/handler"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var GetSecretCmd = &cobra.Command{
	Use:     "secrets [name]",
	Short:   "Get secrets.",
	Long:    ``,
	Aliases: []string{"secret"},
	Args:    cobra.RangeArgs(0, 1),

	Run: func(cmd *cobra.Command, args []string) {
		c := client.NewSecretV1BetaApi()

		if len(args) == 0 {
			secretList, err := c.ListSecrets()

			utils.Check(err)

			const padding = 3
			w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)

			fmt.Fprintln(w, "NAME\tAGE")

			for _, s := range secretList.Secrets {
				updateTime, err := s.Metadata.UpdateTime.Int64()

				utils.Check(err)

				fmt.Fprintf(w, "%s\t%s\n", s.Metadata.Name, handler.RelativeAgeForHumans(updateTime))
			}

			w.Flush()
		} else {
			name := args[0]

			secret, err := c.GetSecret(name)

			utils.Check(err)

			y, err := secret.ToYaml()

			utils.Check(err)

			fmt.Printf("%s", y)
		}
	},
}
