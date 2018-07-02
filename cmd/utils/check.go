package utils

import (
	"fmt"
	"log"
	"os"
)

func Check(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", message)

		log.Fatal("error: %+v\n", err)
	}
}

func Fail(message string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", message)

	os.Exit(1)
}
