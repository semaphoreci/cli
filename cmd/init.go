package cmd

import (
	"fmt"
	"os"

	"github.com/renderedtext/sem/client"
	"github.com/renderedtext/sem/cmd/utils"
	"github.com/renderedtext/sem/config"
	"github.com/renderedtext/sem/generators"

	"github.com/spf13/cobra"
	"github.com/tcnksm/go-gitconfig"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a project",
	Long:  ``,

	Run: func(cmd *cobra.Command, args []string) {
		RunInit(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func RunInit(cmd *cobra.Command, args []string) {
	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		utils.Fail("not a git repository")
	}

	if _, err := os.Stat(".semaphore/semaphore.yml"); err == nil {
		utils.Fail(".semaphore/semaphore.yml is already present in the repository")
	}

	repo_url, err := gitconfig.OriginURL()

	utils.CheckWithMessage(err, "failed to extract remote from git configuration")

	project_yaml, err := generators.GenerateProjectYamlFromRepoUrl(repo_url)

	utils.Check(err)

	project, err := client.InitProjectFromYaml(project_yaml)

	utils.Check(err)

	err = generators.GeneratePipelineYaml()

	utils.Check(err)

	fmt.Printf("Project is created. You can find it at https://%s/projects/%s.\n", config.GetHost(), project.Metadata.Name)
	fmt.Println("")
	fmt.Printf("To run your first pipeline execute:\n")
	fmt.Println("")
	fmt.Printf("  git add .semaphore/semaphore.yml && git commit -m \"First pipeline\" && git push\n")
	fmt.Println("")
}

func registerProjectOnSemaphore(project_yaml []byte) (string, error) {

}
