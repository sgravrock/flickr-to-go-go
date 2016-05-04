package dl

import (
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

type Downloader interface {
	DownloadPhotolist(flickr flickrapi.Client, fs storage.Storage) ([]flickrapi.PhotoListEntry, error)
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

	err = fs.WriteJson("photolist.json", photos)
	if err != nil {
		return nil, err
	}

	return photos, nil
}
