package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-openapi/runtime"
)

//
// Checks if an error is present.
//
// If it is present, it displays the error and exits with status 1.
//
// If you want to display a custom message use CheckWithMessage.
//
func Check(err error) {
	if apiError, success := err.(*runtime.APIError); success {
		// we are dealing with API error

		response, success := apiError.Response.(http.Response)

		if success {
			body, err := ioutil.ReadAll(response.Body)

			if err == nil {
				fmt.Fprintf(os.Stderr, "error: (status %d) %#v\n", apiError.Code, body)
			}
		} else {
			fmt.Fprintf(os.Stderr, "error: (status %d) %#v\n", apiError.Code, apiError.Response)
		}

		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err.Error())

		os.Exit(1)
	}
}

//
// Checks if an error is present.
//
// If it is present, it displays the provided message and exits with status 1.
//
func CheckWithMessage(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %+v\n", message)
	}
}

func Fail(message string) {
	fmt.Fprintf(os.Stderr, "error: %s\n", message)

	os.Exit(1)
}
