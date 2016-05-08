package dl

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

type Downloader interface {
	DownloadPhotolist(flickr flickrapi.Client, fs storage.Storage) ([]flickrapi.PhotoListEntry, error)
	DownloadPhotoInfo(flickr flickrapi.Client, fs storage.Storage, id string) error
	DownloadOriginal(httpClient *http.Client, fs storage.Storage,
		photo flickrapi.PhotoListEntry) error
}

func NewDownloader(stdout io.Writer) Downloader {
	return &downloader{stdout}
}

type downloader struct {
	stdout io.Writer
}

func (d *downloader) DownloadPhotolist(client flickrapi.Client,
	fs storage.Storage) ([]flickrapi.PhotoListEntry, error) {

	fmt.Fprintln(d.stdout, "Downloading photo list")
	photos, err := client.GetPhotos(500)
	if err != nil {
		return nil, err
	}

	err = savePhotolist(fs, photos)
	if err != nil {
		return nil, err
	}

	return photos, nil
}

func savePhotolist(fs storage.Storage, photos []flickrapi.PhotoListEntry) error {
	toSave := make([]map[string]interface{}, len(photos))

	for i, p := range photos {
		toSave[i] = p.Data
	}

	return fs.WriteJson("photolist.json", toSave)
}

func (dl *downloader) DownloadPhotoInfo(flickr flickrapi.Client,
	fs storage.Storage, id string) error {

	path := fmt.Sprintf("photo-info/%s.json", id)
	if fs.Exists(path) {
		return nil
	}

	fmt.Fprintf(dl.stdout, "Downloading info for photo %s\n", id)
	info, err := flickr.GetPhotoInfo(id)
	if err != nil {
		return err
	}

	return fs.WriteJson(path, info)
}

func (dl *downloader) DownloadOriginal(httpClient *http.Client,
	fs storage.Storage, photo flickrapi.PhotoListEntry) error {
	id, err := photo.Id()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("originals/%s.jpg", id)
	if fs.Exists(path) {
		return nil
	}

	fmt.Fprintf(dl.stdout, "Downloading original of photo %s\n", id)
	url, err := photo.OriginalUrl()
	if err != nil {
		return err
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("Request for %s returned status %d",
			url, resp.StatusCode)
		return errors.New(msg)
	}

	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f, err := fs.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(contents)
	return err
}
