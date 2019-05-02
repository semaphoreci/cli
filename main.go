package main

import "github.com/semaphoreci/cli/cmd"

// adjusted to real version during building
var VERSION = "dev"

func main() {
	cmd.VERSION = VERSION

	cmd.Execute()
}
