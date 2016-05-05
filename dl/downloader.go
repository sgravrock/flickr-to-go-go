package dl

import (
	"fmt"

	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

type Downloader interface {
	DownloadPhotolist(flickr flickrapi.Client, fs storage.Storage) ([]flickrapi.PhotoListEntry, error)
	DownloadPhotoInfo(flickr flickrapi.Client, fs storage.Storage, id string) error
}

func NewDownloader() Downloader {
	return &downloader{}
}

type downloader struct{}

func (d *downloader) DownloadPhotolist(client flickrapi.Client,
	fs storage.Storage) ([]flickrapi.PhotoListEntry, error) {

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

	info, err := flickr.GetPhotoInfo(id)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("photo-info/%s.json", id)
	return fs.WriteJson(path, info)
}
