// Code generated by counterfeiter. DO NOT EDIT.
package flickrapifakes

import (
	"sync"

	"github.com/sgravrock/flickr-to-go-go/flickrapi"
)

type FakeClient struct {
	GetStub        func(method string, params map[string]string) (map[string]interface{}, error)
	getMutex       sync.RWMutex
	getArgsForCall []struct {
		method string
		params map[string]string
	}
	getReturns struct {
		result1 map[string]interface{}
		result2 error
	}
	getReturnsOnCall map[int]struct {
		result1 map[string]interface{}
		result2 error
	}
	GetUsernameStub        func() (string, error)
	getUsernameMutex       sync.RWMutex
	getUsernameArgsForCall []struct{}
	getUsernameReturns     struct {
		result1 string
		result2 error
	}
	getUsernameReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetRecentPhotoIdsStub        func(timestamp uint32, pageSize int) ([]string, error)
	getRecentPhotoIdsMutex       sync.RWMutex
	getRecentPhotoIdsArgsForCall []struct {
		timestamp uint32
		pageSize  int
	}
	getRecentPhotoIdsReturns struct {
		result1 []string
		result2 error
	}
	getRecentPhotoIdsReturnsOnCall map[int]struct {
		result1 []string
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
	getPhotosReturnsOnCall map[int]struct {
		result1 []flickrapi.PhotoListEntry
		result2 error
	}
	GetPhotoInfoStub        func(photoId string) (map[string]interface{}, error)
	getPhotoInfoMutex       sync.RWMutex
	getPhotoInfoArgsForCall []struct {
		photoId string
	}
	getPhotoInfoReturns struct {
		result1 map[string]interface{}
		result2 error
	}
	getPhotoInfoReturnsOnCall map[int]struct {
		result1 map[string]interface{}
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeClient) Get(method string, params map[string]string) (map[string]interface{}, error) {
	fake.getMutex.Lock()
	ret, specificReturn := fake.getReturnsOnCall[len(fake.getArgsForCall)]
	fake.getArgsForCall = append(fake.getArgsForCall, struct {
		method string
		params map[string]string
	}{method, params})
	fake.recordInvocation("Get", []interface{}{method, params})
	fake.getMutex.Unlock()
	if fake.GetStub != nil {
		return fake.GetStub(method, params)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getReturns.result1, fake.getReturns.result2
}

func (fake *FakeClient) GetCallCount() int {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return len(fake.getArgsForCall)
}

func (fake *FakeClient) GetArgsForCall(i int) (string, map[string]string) {
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	return fake.getArgsForCall[i].method, fake.getArgsForCall[i].params
}

func (fake *FakeClient) GetReturns(result1 map[string]interface{}, result2 error) {
	fake.GetStub = nil
	fake.getReturns = struct {
		result1 map[string]interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetReturnsOnCall(i int, result1 map[string]interface{}, result2 error) {
	fake.GetStub = nil
	if fake.getReturnsOnCall == nil {
		fake.getReturnsOnCall = make(map[int]struct {
			result1 map[string]interface{}
			result2 error
		})
	}
	fake.getReturnsOnCall[i] = struct {
		result1 map[string]interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetUsername() (string, error) {
	fake.getUsernameMutex.Lock()
	ret, specificReturn := fake.getUsernameReturnsOnCall[len(fake.getUsernameArgsForCall)]
	fake.getUsernameArgsForCall = append(fake.getUsernameArgsForCall, struct{}{})
	fake.recordInvocation("GetUsername", []interface{}{})
	fake.getUsernameMutex.Unlock()
	if fake.GetUsernameStub != nil {
		return fake.GetUsernameStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getUsernameReturns.result1, fake.getUsernameReturns.result2
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

func (fake *FakeClient) GetUsernameReturnsOnCall(i int, result1 string, result2 error) {
	fake.GetUsernameStub = nil
	if fake.getUsernameReturnsOnCall == nil {
		fake.getUsernameReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getUsernameReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetRecentPhotoIds(timestamp uint32, pageSize int) ([]string, error) {
	fake.getRecentPhotoIdsMutex.Lock()
	ret, specificReturn := fake.getRecentPhotoIdsReturnsOnCall[len(fake.getRecentPhotoIdsArgsForCall)]
	fake.getRecentPhotoIdsArgsForCall = append(fake.getRecentPhotoIdsArgsForCall, struct {
		timestamp uint32
		pageSize  int
	}{timestamp, pageSize})
	fake.recordInvocation("GetRecentPhotoIds", []interface{}{timestamp, pageSize})
	fake.getRecentPhotoIdsMutex.Unlock()
	if fake.GetRecentPhotoIdsStub != nil {
		return fake.GetRecentPhotoIdsStub(timestamp, pageSize)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getRecentPhotoIdsReturns.result1, fake.getRecentPhotoIdsReturns.result2
}

func (fake *FakeClient) GetRecentPhotoIdsCallCount() int {
	fake.getRecentPhotoIdsMutex.RLock()
	defer fake.getRecentPhotoIdsMutex.RUnlock()
	return len(fake.getRecentPhotoIdsArgsForCall)
}

func (fake *FakeClient) GetRecentPhotoIdsArgsForCall(i int) (uint32, int) {
	fake.getRecentPhotoIdsMutex.RLock()
	defer fake.getRecentPhotoIdsMutex.RUnlock()
	return fake.getRecentPhotoIdsArgsForCall[i].timestamp, fake.getRecentPhotoIdsArgsForCall[i].pageSize
}

func (fake *FakeClient) GetRecentPhotoIdsReturns(result1 []string, result2 error) {
	fake.GetRecentPhotoIdsStub = nil
	fake.getRecentPhotoIdsReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetRecentPhotoIdsReturnsOnCall(i int, result1 []string, result2 error) {
	fake.GetRecentPhotoIdsStub = nil
	if fake.getRecentPhotoIdsReturnsOnCall == nil {
		fake.getRecentPhotoIdsReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.getRecentPhotoIdsReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetPhotos(pageSize int) ([]flickrapi.PhotoListEntry, error) {
	fake.getPhotosMutex.Lock()
	ret, specificReturn := fake.getPhotosReturnsOnCall[len(fake.getPhotosArgsForCall)]
	fake.getPhotosArgsForCall = append(fake.getPhotosArgsForCall, struct {
		pageSize int
	}{pageSize})
	fake.recordInvocation("GetPhotos", []interface{}{pageSize})
	fake.getPhotosMutex.Unlock()
	if fake.GetPhotosStub != nil {
		return fake.GetPhotosStub(pageSize)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPhotosReturns.result1, fake.getPhotosReturns.result2
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

func (fake *FakeClient) GetPhotosReturnsOnCall(i int, result1 []flickrapi.PhotoListEntry, result2 error) {
	fake.GetPhotosStub = nil
	if fake.getPhotosReturnsOnCall == nil {
		fake.getPhotosReturnsOnCall = make(map[int]struct {
			result1 []flickrapi.PhotoListEntry
			result2 error
		})
	}
	fake.getPhotosReturnsOnCall[i] = struct {
		result1 []flickrapi.PhotoListEntry
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetPhotoInfo(photoId string) (map[string]interface{}, error) {
	fake.getPhotoInfoMutex.Lock()
	ret, specificReturn := fake.getPhotoInfoReturnsOnCall[len(fake.getPhotoInfoArgsForCall)]
	fake.getPhotoInfoArgsForCall = append(fake.getPhotoInfoArgsForCall, struct {
		photoId string
	}{photoId})
	fake.recordInvocation("GetPhotoInfo", []interface{}{photoId})
	fake.getPhotoInfoMutex.Unlock()
	if fake.GetPhotoInfoStub != nil {
		return fake.GetPhotoInfoStub(photoId)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPhotoInfoReturns.result1, fake.getPhotoInfoReturns.result2
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

func (fake *FakeClient) GetPhotoInfoReturns(result1 map[string]interface{}, result2 error) {
	fake.GetPhotoInfoStub = nil
	fake.getPhotoInfoReturns = struct {
		result1 map[string]interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) GetPhotoInfoReturnsOnCall(i int, result1 map[string]interface{}, result2 error) {
	fake.GetPhotoInfoStub = nil
	if fake.getPhotoInfoReturnsOnCall == nil {
		fake.getPhotoInfoReturnsOnCall = make(map[int]struct {
			result1 map[string]interface{}
			result2 error
		})
	}
	fake.getPhotoInfoReturnsOnCall[i] = struct {
		result1 map[string]interface{}
		result2 error
	}{result1, result2}
}

func (fake *FakeClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getMutex.RLock()
	defer fake.getMutex.RUnlock()
	fake.getUsernameMutex.RLock()
	defer fake.getUsernameMutex.RUnlock()
	fake.getRecentPhotoIdsMutex.RLock()
	defer fake.getRecentPhotoIdsMutex.RUnlock()
	fake.getPhotosMutex.RLock()
	defer fake.getPhotosMutex.RUnlock()
	fake.getPhotoInfoMutex.RLock()
	defer fake.getPhotoInfoMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeClient) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ flickrapi.Client = new(FakeClient)
