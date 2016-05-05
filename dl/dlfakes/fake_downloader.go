// This file was generated by counterfeiter
package dlfakes

import (
	"net/http"
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
	DownloadPhotoInfoStub        func(flickr flickrapi.Client, fs storage.Storage, id string) error
	downloadPhotoInfoMutex       sync.RWMutex
	downloadPhotoInfoArgsForCall []struct {
		flickr flickrapi.Client
		fs     storage.Storage
		id     string
	}
	downloadPhotoInfoReturns struct {
		result1 error
	}
	DownloadOriginalStub        func(httpClient *http.Client, fs storage.Storage, photo flickrapi.PhotoListEntry) error
	downloadOriginalMutex       sync.RWMutex
	downloadOriginalArgsForCall []struct {
		httpClient *http.Client
		fs         storage.Storage
		photo      flickrapi.PhotoListEntry
	}
	downloadOriginalReturns struct {
		result1 error
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

func (fake *FakeDownloader) DownloadPhotoInfo(flickr flickrapi.Client, fs storage.Storage, id string) error {
	fake.downloadPhotoInfoMutex.Lock()
	fake.downloadPhotoInfoArgsForCall = append(fake.downloadPhotoInfoArgsForCall, struct {
		flickr flickrapi.Client
		fs     storage.Storage
		id     string
	}{flickr, fs, id})
	fake.downloadPhotoInfoMutex.Unlock()
	if fake.DownloadPhotoInfoStub != nil {
		return fake.DownloadPhotoInfoStub(flickr, fs, id)
	} else {
		return fake.downloadPhotoInfoReturns.result1
	}
}

func (fake *FakeDownloader) DownloadPhotoInfoCallCount() int {
	fake.downloadPhotoInfoMutex.RLock()
	defer fake.downloadPhotoInfoMutex.RUnlock()
	return len(fake.downloadPhotoInfoArgsForCall)
}

func (fake *FakeDownloader) DownloadPhotoInfoArgsForCall(i int) (flickrapi.Client, storage.Storage, string) {
	fake.downloadPhotoInfoMutex.RLock()
	defer fake.downloadPhotoInfoMutex.RUnlock()
	return fake.downloadPhotoInfoArgsForCall[i].flickr, fake.downloadPhotoInfoArgsForCall[i].fs, fake.downloadPhotoInfoArgsForCall[i].id
}

func (fake *FakeDownloader) DownloadPhotoInfoReturns(result1 error) {
	fake.DownloadPhotoInfoStub = nil
	fake.downloadPhotoInfoReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeDownloader) DownloadOriginal(httpClient *http.Client, fs storage.Storage, photo flickrapi.PhotoListEntry) error {
	fake.downloadOriginalMutex.Lock()
	fake.downloadOriginalArgsForCall = append(fake.downloadOriginalArgsForCall, struct {
		httpClient *http.Client
		fs         storage.Storage
		photo      flickrapi.PhotoListEntry
	}{httpClient, fs, photo})
	fake.downloadOriginalMutex.Unlock()
	if fake.DownloadOriginalStub != nil {
		return fake.DownloadOriginalStub(httpClient, fs, photo)
	} else {
		return fake.downloadOriginalReturns.result1
	}
}

func (fake *FakeDownloader) DownloadOriginalCallCount() int {
	fake.downloadOriginalMutex.RLock()
	defer fake.downloadOriginalMutex.RUnlock()
	return len(fake.downloadOriginalArgsForCall)
}

func (fake *FakeDownloader) DownloadOriginalArgsForCall(i int) (*http.Client, storage.Storage, flickrapi.PhotoListEntry) {
	fake.downloadOriginalMutex.RLock()
	defer fake.downloadOriginalMutex.RUnlock()
	return fake.downloadOriginalArgsForCall[i].httpClient, fake.downloadOriginalArgsForCall[i].fs, fake.downloadOriginalArgsForCall[i].photo
}

func (fake *FakeDownloader) DownloadOriginalReturns(result1 error) {
	fake.DownloadOriginalStub = nil
	fake.downloadOriginalReturns = struct {
		result1 error
	}{result1}
}

var _ dl.Downloader = new(FakeDownloader)
