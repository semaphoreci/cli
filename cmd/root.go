package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/fatih/color"
	"github.com/semaphoreci/cli/cmd/utils"
	"github.com/semaphoreci/cli/config"

	versions "github.com/hashicorp/go-version"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Verbose bool
var NoVersionUpdateChecks bool

var versionCheckInterval = time.Minute
var installationURL string = "https://storage.googleapis.com/sem-cli-releases/get.sh"
var installationVersionRegex string = `VERSION="(.*)"`
var newVersionWarning = fmt.Sprintf(
	`
A new version of the Semaphore CLI (%%s) is available. You can install it with:
  'curl %s | bash'
`, installationURL)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "sem",
	Short: "Semaphore 2.0 command line interface",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !Verbose {
			log.SetOutput(ioutil.Discard)
		}
	},

	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute(version string) {
	cmdErr := RootCmd.Execute()
	if cmdErr != nil {
		fmt.Println(cmdErr)
	}

	//
	// Before exiting, check for newer versions of this CLI periodically.
	// If an error happens looking for the latest version, just let it go.
	//
	if !NoVersionUpdateChecks && shouldCheckVersion() {
		v, err := findLatestVersion()
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}

		currentVersion, _ := versions.NewVersion(version)
		latestVersion, _ := versions.NewVersion(v)
		if latestVersion.GreaterThan(currentVersion) {
			color.Set(color.FgYellow)
			fmt.Printf(newVersionWarning, latestVersion)
			color.Unset()
		}
	}

	//
	// Now, we can exit with the appropriate exit code.
	//
	if cmdErr != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	RootCmd.PersistentFlags().BoolVarP(&NoVersionUpdateChecks, "no-new-version-checks", "", false, "Do not check for new versions of this CLI")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	home, err := homedir.Dir()

	utils.CheckWithMessage(err, "failed to find home directory")

	// Search config in home directory with name ".sem" (without extension).
	viper.AddConfigPath(home)
	viper.SetConfigName(".sem")

	// Touch config file and make sure that it exists
	path := fmt.Sprintf("%s/.sem.yaml", home)

	// #nosec
	_, _ = os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)

	err = viper.ReadInConfig()

	utils.CheckWithMessage(err, "failed to load config file")
}

func shouldCheckVersion() bool {
	now := time.Now().UTC().Format(time.RFC3339)
	t := config.GetVersionCheckTime()

	//
	// If the config is empty,
	// it means we never checked the version.
	if t == "" {
		//
		// TODO: this fails when the CLI is run without a command,
		// because initConfig only executes after an actual command is run.
		//
		config.SetVersionCheckTime(now)
		return true
	}

	lastCheck, err := time.Parse(time.RFC3339, t)

	//
	// If we have a config value, but we can't parse,
	// something really wrong happened, so just reset it.
	//
	if err != nil {
		config.SetVersionCheckTime(now)
		return true
	}

	if time.Since(lastCheck) > versionCheckInterval {

		//
		// TODO: this fails when the CLI is run without a command,
		// because initConfig only executes after an actual command is run.
		//
		config.SetVersionCheckTime(now)
		return true
	}

	return false
}

func findLatestVersion() (string, error) {
	response, err := http.Get(installationURL)
	if err != nil {
		return "", err
	}

	script, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	r, err := regexp.Compile(installationVersionRegex)
	if err != nil {
		return "", err
	}

	matches := r.FindStringSubmatch(string(script))
	if len(matches) < 2 {
		return "", fmt.Errorf("could not find newer version in installation script")
	}

	return matches[1], nil
}
