// This file was generated by counterfeiter
package flickrapifakes

import (
	"sync"

	"github.com/sgravrock/flickr-to-go-go/flickrapi"
)

type FakeClient struct {
	GetStub        func(method string, params map[string]string, payload flickrapi.FlickrPayload) error
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		method  string
		params  map[string]string
		payload flickrapi.FlickrPayload
	}
	getReturns struct {
		result1 error
	}
	GetUsernameStub        func() (string, error)
	getUsernameMutex       sync.RWMutex
	getUsernameArgsForCall []struct{}
	getUsernameReturns     struct {
		result1 string
		result2 error
	}
	GetPhotosStub        func(pageSize int) ([]flickrapi.PhotoListEntry, error)
	getPhotosMutex       sync.RWMutex
	getPhotosArgsForCall []struct {
		pageSize int
	}
	getPhotosReturns struct {
		result1 []flickrapi.PhotoListEntry
		result2 error
	}
	GetPhotoInfoStub        func(photoId string) (flickrapi.PhotoInfo, error)
	getPhotoInfoMutex       sync.RWMutex
	getPhotoInfoArgsForCall []struct {
		photoId string
	}
	getPhotoInfoReturns struct {
		result1 flickrapi.PhotoInfo
		result2 error
	}
}

func (fake *FakeClient) Get(method string, params map[string]string, payload flickrapi.FlickrPayload) error {
	fake.getMutex.Lock()
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		method  string
		params  map[string]string
		payload flickrapi.FlickrPayload
	}{method, params, payload})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(method, params, payload)
	} else {
		return fake.getReturns.result1
	}
}

func (fake *FakeClient) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeClient) GetArgsForCall(i int) (string, map[string]string, flickrapi.FlickrPayload) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].method, fake.getArgsForCall[i].params, fake.getArgsForCall[i].payload
}

func (fake *FakeClient) GetReturns(result1 error) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeClient) GetUsername() (string, error) {
	fake.getUsernameMutex.Lock()
	fake.getUsernameArgsForCall = append(fake.getUsernameArgsForCall, struct{}{})
	fake.getUsernameMutex.Unlock()
	if fake.GetUsernameStub != nil {
		return fake.GetUsernameStub()
	} else {
		return fake.getUsernameReturns.result1, fake.getUsernameReturns.result2
	}
}

func (fake *FakeClient) GetUsernameCallCount() int {
	fake.getUsernameMutex.RLock()
	defer fake.getUsernameMutex.RUnlock()
	return len(fake.getUsernameArgsForCall)
}

func (fake *FakeClient) GetUsernameReturns(result1 string, result2 error) {
	fake.GetUsernameStub = nil
	fake.getUsernameReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetPhotos(pageSize int) ([]flickrapi.PhotoListEntry, error) {
	fake.getPhotosMutex.Lock()
	fake.getPhotosArgsForCall = append(fake.getPhotosArgsForCall, struct {
		pageSize int
	}{pageSize})
	fake.getPhotosMutex.Unlock()
	if fake.GetPhotosStub != nil {
		return fake.GetPhotosStub(pageSize)
	} else {
		return fake.getPhotosReturns.result1, fake.getPhotosReturns.result2
	}
}

func (fake *FakeClient) GetPhotosCallCount() int {
	fake.getPhotosMutex.RLock()
	defer fake.getPhotosMutex.RUnlock()
	return len(fake.getPhotosArgsForCall)
}

func (fake *FakeClient) GetPhotosArgsForCall(i int) int {
	fake.getPhotosMutex.RLock()
	defer fake.getPhotosMutex.RUnlock()
	return fake.getPhotosArgsForCall[i].pageSize
}

func (fake *FakeClient) GetPhotosReturns(result1 []flickrapi.PhotoListEntry, result2 error) {
	fake.GetPhotosStub = nil
	fake.getPhotosReturns = struct {
		result1 []flickrapi.PhotoListEntry
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetPhotoInfo(photoId string) (flickrapi.PhotoInfo, error) {
	fake.getPhotoInfoMutex.Lock()
	fake.getPhotoInfoArgsForCall = append(fake.getPhotoInfoArgsForCall, struct {
		photoId string
	}{photoId})
	fake.getPhotoInfoMutex.Unlock()
	if fake.GetPhotoInfoStub != nil {
		return fake.GetPhotoInfoStub(photoId)
	} else {
		return fake.getPhotoInfoReturns.result1, fake.getPhotoInfoReturns.result2
	}
}

func (fake *FakeClient) GetPhotoInfoCallCount() int {
	fake.getPhotoInfoMutex.RLock()
	defer fake.getPhotoInfoMutex.RUnlock()
	return len(fake.getPhotoInfoArgsForCall)
}

func (fake *FakeClient) GetPhotoInfoArgsForCall(i int) string {
	fake.getPhotoInfoMutex.RLock()
	defer fake.getPhotoInfoMutex.RUnlock()
	return fake.getPhotoInfoArgsForCall[i].photoId
}

func (fake *FakeClient) GetPhotoInfoReturns(result1 flickrapi.PhotoInfo, result2 error) {
	fake.GetPhotoInfoStub = nil
	fake.getPhotoInfoReturns = struct {
		result1 flickrapi.PhotoInfo
		result2 error
	}{result1, result2}
}

var _ flickrapi.Client = new(FakeClient)
