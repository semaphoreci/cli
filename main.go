package main

import "github.com/semaphoreci/cli/cmd"

// injected as ldflags during building
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// inject version information
	cmd.ReleaseVersion = version
	cmd.ReleaseCommit = commit
	cmd.ReleaseDate = date

	cmd.Execute()
}
