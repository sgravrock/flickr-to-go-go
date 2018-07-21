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

	ftg := application{
		baseUrl:       baseUrl,
		savecreds:     savecreds,
		authenticator: authenticator,
		downloader:    downloader,
		fileStore:     fileStore,
		clock:         clock,
		stdout:        stdout,
		stderr:        stderr,
	}
	return ftg.Run()
}

type application struct {
	baseUrl       string
	savecreds     bool
	authenticator auth.Authenticator
	downloader    dl.Downloader
	fileStore     storage.Storage
	clock         clock.Clock
	stdout        io.Writer
	stderr        io.Writer
}

func (ftg *application) Run() int {
	err := ftg.fileStore.EnsureRoot()
	if err != nil {
		msg := err.Error()[5:] // strip leading "stat "
		fmt.Fprintln(ftg.stderr, msg)
		return 1
	}

	httpClient, err := ftg.authenticator.Authenticate(ftg.savecreds)
	if err != nil {
		fmt.Fprintln(ftg.stderr, err.Error())
		return 1
	}
	flickrClient := flickrapi.NewClient(httpClient, ftg.baseUrl)

	photos, err := ftg.downloader.DownloadPhotolist(flickrClient, ftg.fileStore)
	if err != nil {
		fmt.Fprintf(ftg.stderr, "Error downloading photo list: %s\n",
			err.Error())
		return 1
	}

	updatedPhotoIds, err := ftg.getUpdatedPhotoIds(flickrClient)
	if err != nil {
		fmt.Fprintf(ftg.stderr, "Error downloading the list of updated photos: %s\n", err.Error())
	}

	fmt.Fprintf(ftg.stdout, "Got %d photos\n", len(photos))

	for _, p := range photos {
		id, err := p.Id()
		if err != nil {
			fmt.Fprintln(ftg.stderr, err.Error())
			return 1
		}

		if ftg.shouldDownloadInfo(id, updatedPhotoIds) {
			err = ftg.downloader.DownloadPhotoInfo(flickrClient,
				ftg.fileStore, id)
			if err != nil {
				fmt.Fprintf(ftg.stderr,
					"Error downloading info for %s: %s\n",
					id, err.Error())
				return 1
			}
		}

		if ftg.shouldDownloadOriginal(id, updatedPhotoIds) {
			err = ftg.downloader.DownloadOriginal(httpClient, ftg.fileStore, p)
			if err != nil {
				fmt.Fprintf(ftg.stderr,
					"Error downloading original for %s: %s\n",
					id, err.Error())
				return 1
			}
		}
	}

	err = ftg.moveDeletedFiles(photos, "photo-info", "info")
	if err != nil {
		fmt.Fprintln(ftg.stderr, err.Error())
		return 1
	}

	err = ftg.moveDeletedFiles(photos, "originals", "original")
	if err != nil {
		fmt.Fprintln(ftg.stderr, err.Error())
		return 1
	}

	err = timestamp.Write(ftg.clock, ftg.fileStore)
	if err != nil {
		fmt.Fprintf(ftg.stderr, "Error saving timestamp: %s\n", err.Error())
		return 1
	}

	return 0
}

func (ftg *application) shouldDownloadOriginal(
	photoId string,
	updatedPhotoIds []string) bool {

	return updatedPhotoIds == nil ||
		containsString(updatedPhotoIds, photoId) ||
		!ftg.downloader.OriginalExists(ftg.fileStore, photoId)
}

func (ftg *application) shouldDownloadInfo(
	photoId string,
	updatedPhotoIds []string) bool {

	return updatedPhotoIds == nil ||
		containsString(updatedPhotoIds, photoId) ||
		!ftg.downloader.PhotoInfoExists(ftg.fileStore, photoId)
}

func (ftg *application) getUpdatedPhotoIds(
	flickr flickrapi.Client) ([]string, error) {

	lastRun := timestamp.Read(ftg.fileStore, ftg.stderr)
	if lastRun == 0 {
		return nil, nil
	}

	return ftg.downloader.GetRecentPhotoIds(lastRun, flickr)
}

func (ftg *application) moveDeletedFiles(
	photos []flickrapi.PhotoListEntry,
	srcDir string,
	description string) error {

	srcFiles, err := ftg.fileStore.ListFiles(srcDir)
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

		if !ftg.containsPhoto(photos, photoId) {
			fmt.Printf("Moving %s of deleted photo %s to attic\n",
				description, photoId)
			oldPath := path.Join(srcDir, filename)
			newPath := path.Join("attic", oldPath)
			err = ftg.fileStore.Move(oldPath, newPath)

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

func (ftg *application) containsPhoto(
	photos []flickrapi.PhotoListEntry,
	idToFind string) bool {

	for _, candidate := range photos {
		id, err := candidate.Id()

		if err != nil {
			// Should've already been caught elsewhere. Warn and continue.
			fmt.Fprintln(ftg.stderr, "Photo has no id")
		} else if id == idToFind {
			return true
		}
	}

	return false
}

func containsString(haystack []string, needle string) bool {
	for _, candidate := range haystack {
		if candidate == needle {
			return true
		}
	}

	return false
}
