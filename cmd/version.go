package cmd

import (
	"fmt"
	"runtime"

	"github.com/hamir-suspect/go-latest"
	"github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
)

var (
	ReleaseVersion = "dev"
	ReleaseCommit  = "none"
	ReleaseDate    = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v, commit %v, built at %v\n",
			ReleaseVersion,
			ReleaseCommit,
			ReleaseDate)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
	if runtime.GOOS == "windows" {
		Reset  = ""
		Yellow = ""
		Green = ""
	}
}

var Reset  = "\033[0m"
var Green= "\033[32m"
var Yellow = "\033[33m"
var outdatedVersionMsgTemplate = "%swarn: You are using outdated version of cli: %s newest version is: %s\nInstall the new version using this snipet: %s%s%s\n"
var InstallScript = "curl https://storage.googleapis.com/sem-cli-releases/get.sh | bash"

func CheckNewerVersion() {
	githubTag := latest.GithubTag{
		Owner: "semaphoreci",
		Repository: "cli",
	}
	
	// ignore version check if using dev version
	_, err := version.NewSemver(ReleaseVersion)
	if err != nil {
		return
	}
	
	res, err := latest.Check(&githubTag, "v"+ReleaseVersion)
	if err!= nil {
		fmt.Println("error checking the version: ", err.Error())
		return
	}
	if res.Outdated {
		fmt.Printf(outdatedVersionMsgTemplate, Yellow, "v"+ReleaseVersion, res.Current, Green, InstallScript, Reset)
	}
}