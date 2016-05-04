// This file was generated by counterfeiter
package dlfakes

import (
	"sync"

	"github.com/sgravrock/flickr-to-go-go/dl"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

type FakeDownloader struct {
	DownloadPhotolistStub        func(flickr flickrapi.Client, fs storage.Storage) ([]flickrapi.PhotoListEntry, error)
	downloadPhotolistMutex       sync.RWMutex
	downloadPhotolistArgsForCall []struct {
		flickr flickrapi.Client
		fs     storage.Storage
	}
	downloadPhotolistReturns struct {
		result1 []flickrapi.PhotoListEntry
		result2 error
	}
}

func (fake *FakeDownloader) DownloadPhotolist(flickr flickrapi.Client, fs storage.Storage) ([]flickrapi.PhotoListEntry, error) {
	fake.downloadPhotolistMutex.Lock()
	fake.downloadPhotolistArgsForCall = append(fake.downloadPhotolistArgsForCall, struct {
		flickr flickrapi.Client
		fs     storage.Storage
	}{flickr, fs})
	fake.downloadPhotolistMutex.Unlock()
	if fake.DownloadPhotolistStub != nil {
		return fake.DownloadPhotolistStub(flickr, fs)
	} else {
		return fake.downloadPhotolistReturns.result1, fake.downloadPhotolistReturns.result2
	}
}

func (fake *FakeDownloader) DownloadPhotolistCallCount() int {
	fake.downloadPhotolistMutex.RLock()
	defer fake.downloadPhotolistMutex.RUnlock()
	return len(fake.downloadPhotolistArgsForCall)
}

func (fake *FakeDownloader) DownloadPhotolistArgsForCall(i int) (flickrapi.Client, storage.Storage) {
	fake.downloadPhotolistMutex.RLock()
	defer fake.downloadPhotolistMutex.RUnlock()
	return fake.downloadPhotolistArgsForCall[i].flickr, fake.downloadPhotolistArgsForCall[i].fs
}

func (fake *FakeDownloader) DownloadPhotolistReturns(result1 []flickrapi.PhotoListEntry, result2 error) {
	fake.DownloadPhotolistStub = nil
	fake.downloadPhotolistReturns = struct {
		result1 []flickrapi.PhotoListEntry
		result2 error
	}{result1, result2}
}

var _ dl.Downloader = new(FakeDownloader)
