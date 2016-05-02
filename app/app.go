package app

import (
	"fmt"
	"io"

	"github.com/sgravrock/flickr-to-go-go/auth"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
)

func Run(baseUrl string, savecreds bool, authenticator auth.Authenticator,
	stdout io.Writer, stderr io.Writer) int {

	httpClient, err := authenticator.Authenticate(savecreds)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	flickrClient := flickrapi.NewClient(httpClient, baseUrl)
	username, err := flickrClient.GetUsername()
	if err != nil {
		fmt.Fprintf(stdout, "Couldn't verify login: %s\n", err.Error())
		return 1
	}

	fmt.Fprintf(stdout, "You are logged in as %s.\n", username)
	return 0
}
