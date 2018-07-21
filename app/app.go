package app

import (
	"fmt"
	"io"
	"path"

	"regexp"

	"errors"

	"os"

	"github.com/sgravrock/flickr-to-go-go/auth"
	"github.com/sgravrock/flickr-to-go-go/clock"
	"github.com/sgravrock/flickr-to-go-go/dl"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage"
	"github.com/sgravrock/flickr-to-go-go/timestamp"
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

		if shouldDownloadInfo(downloader, fileStore, id, updatedPhotoIds) {
			err = downloader.DownloadPhotoInfo(flickrClient, fileStore, id)
			if err != nil {
				fmt.Fprintf(stderr, "Error downloading info for %s: %s\n",
					id, err.Error())
				return 1
			}
		}

		if shouldDownloadOriginal(downloader, fileStore, id, updatedPhotoIds) {
			err = downloader.DownloadOriginal(httpClient, fileStore, p)
			if err != nil {
				fmt.Fprintf(stderr, "Error downloading original for %s: %s\n",
					id, err.Error())
				return 1
			}
		}
	}

	err = moveDeletedFiles(fileStore, stderr, photos, "photo-info",
		"info")
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	err = moveDeletedFiles(fileStore, stderr, photos, "originals",
		"original")
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	err = timestamp.Write(clock, fileStore)
	if err != nil {
		fmt.Fprintf(stderr, "Error saving timestamp: %s\n", err.Error())
		return 1
	}

	return 0
}

func shouldDownloadOriginal(downloader dl.Downloader,
	fileStore storage.Storage,
	photoId string,
	updatedPhotoIds []string) bool {

	return updatedPhotoIds == nil ||
		containsString(updatedPhotoIds, photoId) ||
		!downloader.OriginalExists(fileStore, photoId)
}

func shouldDownloadInfo(downloader dl.Downloader,
	fileStore storage.Storage,
	photoId string,
	updatedPhotoIds []string) bool {

	return updatedPhotoIds == nil ||
		containsString(updatedPhotoIds, photoId) ||
		!downloader.PhotoInfoExists(fileStore, photoId)
}

func containsString(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if candidate == needle {
			return true
		}
	}

	return false
}

func getUpdatedPhotoIds(flickr flickrapi.Client, downloader dl.Downloader,
	fileStore storage.Storage, stderr io.Writer) ([]string, error) {

	lastRun := timestamp.Read(fileStore, stderr)
	if lastRun == 0 {
		return nil, nil
	}

	return downloader.GetRecentPhotoIds(lastRun, flickr)
}

func moveDeletedFiles(fileStore storage.Storage,
	stderr io.Writer,
	photos []flickrapi.PhotoListEntry,
	srcDir string,
	description string) error {

	srcFiles, err := fileStore.ListFiles(srcDir)
	if err != nil {
		if os.IsNotExist(err) {
			// First run. The directory doesn't exist yet, so there can be
			// nothing to move from it.
			return nil
		}

		msg := fmt.Sprintf("Error reading dir %s: %s\n",
			srcDir, err.Error())
		return errors.New(msg)
	}

	stripExtension := regexp.MustCompile("\\.[^\\.]+$")

	for _, filename := range srcFiles {
		photoId := stripExtension.ReplaceAllLiteralString(filename, "")

		if !containsPhoto(photos, photoId, stderr) {
			fmt.Printf("Moving %s of deleted photo %s to attic\n",
				description, photoId)
			oldPath := path.Join(srcDir, filename)
			newPath := path.Join("attic", oldPath)
			err = fileStore.Move(oldPath, newPath)

			if err != nil {
				msg := fmt.Sprintf(
					"Error moving %s of deleted photo %s: %s\n",
					description, photoId, err.Error())
				return errors.New(msg)
			}
		}
	}

	return nil
}

func containsPhoto(photos []flickrapi.PhotoListEntry,
	idToFind string,
	stderr io.Writer) bool {

	for _, candidate := range photos {
		id, err := candidate.Id()

		if err != nil {
			// Should've already been caught elsewhere. Warn and continue.
			fmt.Fprintln(stderr, "Photo has no id")
		} else if id == idToFind {
			return true
		}
	}

	return false
}
