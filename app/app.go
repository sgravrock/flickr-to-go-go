package app

import (
	"fmt"
	"io"
	"strings"

	"github.com/sgravrock/flickr-to-go-go/auth"
	"github.com/sgravrock/flickr-to-go-go/clock"
	"github.com/sgravrock/flickr-to-go-go/dl"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage"
	"path"
)

func Run(baseUrl string, savecreds bool, authenticator auth.Authenticator,
	downloader dl.Downloader, fileStore storage.Storage, clock clock.Clock,
	stdout io.Writer, stderr io.Writer) int {

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

	updatedPhotoIds, err := getUpdatedPhotoIds(flickrClient, downloader, fileStore, stderr)
	if err != nil {
		fmt.Fprintf(stderr, "Error downloading the list of updated photos: %s\n", err.Error())
	}

	fmt.Fprintf(stdout, "Got %d photos\n", len(photos))

	for _, p := range photos {
		id, err := p.Id()
		if err != nil {
			fmt.Fprintln(stderr, err.Error())
			return 1
		}

		if (shouldDownloadInfo(downloader, fileStore, id, updatedPhotoIds)) {
			err = downloader.DownloadPhotoInfo(flickrClient, fileStore, id)
			if err != nil {
				fmt.Fprintf(stderr, "Error downloading info for %s: %s\n",
					id, err.Error())
				return 1
			}
		}

		if (shouldDownloadOriginal(downloader, fileStore, id, updatedPhotoIds)) {
			err = downloader.DownloadOriginal(httpClient, fileStore, p)
			if err != nil {
				fmt.Fprintf(stderr, "Error downloading original for %s: %s\n",
					id, err.Error())
				return 1
			}
		}

		err = writeTimestamp(clock, fileStore)
		if err != nil {
			fmt.Fprintf(stderr, "Error saving timestamp: %s\n", err.Error())
			return 1
		}
	}

	files, err := fileStore.ListFiles("photo-info")
	if err != nil {
		fmt.Fprintf(stderr, "Error reading info dir: %s\n", err.Error())
		return 1
	}

	for _, filename := range files {
		photoId := strings.Replace(filename, ".json", "", 1)

		if !containsString(updatedPhotoIds, photoId) {
			oldPath := path.Join("photo-info", filename)
			newPath := path.Join("attic", oldPath)
			fileStore.Move(oldPath, newPath)
		}
	}

	return 0
}

func shouldDownloadOriginal(downloader dl.Downloader,
	fileStore storage.Storage,
	photoId string,
	updatedPhotoIds []string) bool {

	return updatedPhotoIds == nil ||
		containsString(updatedPhotoIds, photoId) ||
		!downloader.OriginalExists(fileStore, photoId);
}

func shouldDownloadInfo(downloader dl.Downloader,
	fileStore storage.Storage,
	photoId string,
	updatedPhotoIds []string) bool {

	return updatedPhotoIds == nil ||
		containsString(updatedPhotoIds, photoId) ||
		!downloader.PhotoInfoExists(fileStore, photoId);
}

func containsString(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if candidate == needle {
			return true
		}
	}

	return false
}

func writeTimestamp(clock clock.Clock, fs storage.Storage) error {
	f, err := fs.Create("timestamp")
	if err != nil {
		return err
	}
	defer f.Close()
	s := fmt.Sprint(clock.Now().Unix()) + "\n"
	f.Write([]byte(s))
	return nil
}

func readTimestamp(fs storage.Storage) (uint32, error) {
	f, err := fs.Open("timestamp")
	if err != nil {
		return 0, err
	}
	defer f.Close()
	var result uint32
	_, err = fmt.Fscan(f, &result)
	return result, err
}

func getUpdatedPhotoIds(flickr flickrapi.Client, downloader dl.Downloader,
	fileStore storage.Storage, stderr io.Writer) ([]string, error) {
	timestamp, err := readTimestamp(fileStore)
	if err != nil {
		fmt.Fprintf(stderr, "Error reading timestamp: %s\n", err.Error())
		timestamp = 0
	}

	if timestamp == 0 {
		return nil, nil
	}

	return downloader.GetRecentPhotoIds(timestamp, flickr)
}
