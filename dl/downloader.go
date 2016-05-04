package dl

import "github.com/sgravrock/flickr-to-go-go/flickrapi"

type Downloader interface {
	DownloadPhotolist(flickr flickrapi.Client) ([]flickrapi.PhotoInfo, error)
}

func NewDownloader() Downloader {
	return &downloader{}
}

type downloader struct{}

func (d *downloader) DownloadPhotolist(client flickrapi.Client) ([]flickrapi.PhotoInfo, error) {
	return client.GetPhotos(500)
}
