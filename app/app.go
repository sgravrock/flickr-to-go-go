package app

import (
	"fmt"
	"io"

	"github.com/sgravrock/flickr-to-go-go/auth"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

func Run(baseUrl string, savecreds bool, authenticator auth.Authenticator,
	fileStore storage.Storage, stdout io.Writer, stderr io.Writer) int {

	err := fileStore.EnsureRoot()
	if err != nil {
		msg := err.Error()[5:] // strip leading "stat "
		fmt.Fprintln(stderr, msg)
		return 1
	}

	httpClient, err := authenticator.Authenticate(savecreds)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	flickrClient := flickrapi.NewClient(httpClient, baseUrl)
	username, err := flickrClient.GetUsername()
	if err != nil {
		fmt.Fprintf(stderr, "Couldn't verify login: %s\n", err.Error())
		return 1
	}

	fmt.Fprintf(stdout, "You are logged in as %s.\n", username)
	fmt.Fprintln(stdout, "Downloading photo list")
	photos, err := flickrClient.GetPhotos(500)
	if err != nil {
		fmt.Fprintf(stderr, "Error downloading phot list: %s\n", err.Error())
		return 1
	}

	fmt.Fprintf(stdout, "Got %d photos\n", len(photos))
	return 0
}
