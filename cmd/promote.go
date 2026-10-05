package cmd

import (
	"fmt"
	"strings"

	client "github.com/semaphoreci/cli/api/client"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/spf13/cobra"
)

var promoteCmd = &cobra.Command{
	Use:   "promote [PIPELINE ID] [PROMOTION NAME]",
	Short: "Trigger a promotion of a pipeline.",
	Long: `Trigger the promotion named PROMOTION NAME on the pipeline PIPELINE ID,
as the "Promote" button does in the UI.

Values for a parameterized promotion are passed with --param NAME=VALUE,
repeated once per parameter. A pipeline that failed or is still running
can only be promoted with --override.`,
	Example: `  sem promote 494b76aa-f3f0-4ecf-b5ef-c389591a01be "Deploy to production"
  sem promote 494b76aa-f3f0-4ecf-b5ef-c389591a01be "Deploy" --param REGION=eu-west-1 --param REPLICAS=3
  sem promote 494b76aa-f3f0-4ecf-b5ef-c389591a01be "Deploy" --override`,
	Args: cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {
		pipelineID := args[0]
		name := args[1]

		rawParams, err := cmd.Flags().GetStringArray("param")
		utils.Check(err)

		parameters, err := parsePromotionParameters(rawParams)
		utils.Check(err)

		override, err := cmd.Flags().GetBool("override")
		utils.Check(err)

		c := client.NewPromotionsV1AlphaApi()
		body, err := c.Trigger(pipelineID, name, override, parameters)
		utils.Check(err)

		if len(strings.TrimSpace(string(body))) > 0 {
			fmt.Printf("%s\n", string(body))
		} else {
			fmt.Printf("Promotion '%s' triggered for pipeline %s.\n", name, pipelineID)
		}
	},
}

func parsePromotionParameters(rawParams []string) (map[string]string, error) {
	parameters := map[string]string{}

	for _, raw := range rawParams {
		parts := strings.SplitN(raw, "=", 2)

		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "" {
			return nil, fmt.Errorf("invalid --param '%s', expected NAME=VALUE", raw)
		}

		parameters[strings.TrimSpace(parts[0])] = parts[1]
	}

	return parameters, nil
}

func init() {
	RootCmd.AddCommand(promoteCmd)

	promoteCmd.Flags().StringArrayP("param", "p", []string{}, "value of a promotion parameter, as NAME=VALUE (repeatable)")
	promoteCmd.Flags().Bool("override", false, "promote even if the pipeline failed or is still running")
}
