package app

import (
	"fmt"
	"io"

	"github.com/sgravrock/flickr-to-go-go/auth"
	"github.com/sgravrock/flickr-to-go-go/dl"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

func Run(baseUrl string, savecreds bool, authenticator auth.Authenticator,
	downloader dl.Downloader, fileStore storage.Storage, stdout io.Writer,
	stderr io.Writer) int {

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
	photos, err := downloader.DownloadPhotolist(flickrClient, fileStore)
	if err != nil {
		fmt.Fprintf(stderr, "Error downloading photo list: %s\n", err.Error())
		return 1
	}

	fmt.Fprintf(stdout, "Got %d photos\n", len(photos))

	for _, p := range photos {
		id, err := p.Id()
		if err != nil {
			fmt.Fprintln(stderr, err.Error())
			return 1
		}

		err = downloader.DownloadPhotoInfo(flickrClient, fileStore, id)
		if err != nil {
			fmt.Fprintf(stderr, "Error downloading info for %s: %s\n",
				id, err.Error())
			return 1
		}

		err = downloader.DownloadOriginal(httpClient, fileStore, p)
		if err != nil {
			fmt.Fprintf(stderr, "Error downloading original for %s: %s\n",
				id, err.Error())
			return 1
		}
	}

	return 0
}
