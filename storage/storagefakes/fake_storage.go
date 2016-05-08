// This file was generated by counterfeiter
package storagefakes

import (
	"sync"

	"github.com/sgravrock/flickr-to-go-go/storage"
)

type FakeStorage struct {
	EnsureRootStub        func() error
	ensureRootMutex       sync.RWMutex
	ensureRootArgsForCall []struct{}
	ensureRootReturns     struct {
		result1 error
	}
	ExistsStub        func(string) bool
	existsMutex       sync.RWMutex
	existsArgsForCall []struct {
		arg1 string
	}
	existsReturns struct {
		result1 bool
	}
	CreateStub        func(name string) (storage.File, error)
	createMutex       sync.RWMutex
	createArgsForCall []struct {
		name string
	}
	createReturns struct {
		result1 storage.File
		result2 error
	}
	OpenStub        func(name string) (storage.File, error)
	openMutex       sync.RWMutex
	openArgsForCall []struct {
		name string
	}
	openReturns struct {
		result1 storage.File
		result2 error
	}
	WriteJsonStub        func(name string, payload interface{}) error
	writeJsonMutex       sync.RWMutex
	writeJsonArgsForCall []struct {
		name    string
		payload interface{}
	}
	writeJsonReturns struct {
		result1 error
	}
}

func (fake *FakeStorage) EnsureRoot() error {
	fake.ensureRootMutex.Lock()
	fake.ensureRootArgsForCall = append(fake.ensureRootArgsForCall, struct{}{})
	fake.ensureRootMutex.Unlock()
	if fake.EnsureRootStub != nil {
		return fake.EnsureRootStub()
	} else {
		return fake.ensureRootReturns.result1
	}
}

func (fake *FakeStorage) EnsureRootCallCount() int {
	fake.ensureRootMutex.RLock()
	defer fake.ensureRootMutex.RUnlock()
	return len(fake.ensureRootArgsForCall)
}

func (fake *FakeStorage) EnsureRootReturns(result1 error) {
	fake.EnsureRootStub = nil
	fake.ensureRootReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeStorage) Exists(arg1 string) bool {
	fake.existsMutex.Lock()
	fake.existsArgsForCall = append(fake.existsArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.existsMutex.Unlock()
	if fake.ExistsStub != nil {
		return fake.ExistsStub(arg1)
	} else {
		return fake.existsReturns.result1
	}
}

func (fake *FakeStorage) ExistsCallCount() int {
	fake.existsMutex.RLock()
	defer fake.existsMutex.RUnlock()
	return len(fake.existsArgsForCall)
}

func (fake *FakeStorage) ExistsArgsForCall(i int) string {
	fake.existsMutex.RLock()
	defer fake.existsMutex.RUnlock()
	return fake.existsArgsForCall[i].arg1
}

func (fake *FakeStorage) ExistsReturns(result1 bool) {
	fake.ExistsStub = nil
	fake.existsReturns = struct {
		result1 bool
	}{result1}
}

func (fake *FakeStorage) Create(name string) (storage.File, error) {
	fake.createMutex.Lock()
	fake.createArgsForCall = append(fake.createArgsForCall, struct {
		name string
	}{name})
	fake.createMutex.Unlock()
	if fake.CreateStub != nil {
		return fake.CreateStub(name)
	} else {
		return fake.createReturns.result1, fake.createReturns.result2
	}
}

func (fake *FakeStorage) CreateCallCount() int {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return len(fake.createArgsForCall)
}

func (fake *FakeStorage) CreateArgsForCall(i int) string {
	fake.createMutex.RLock()
	defer fake.createMutex.RUnlock()
	return fake.createArgsForCall[i].name
}

func (fake *FakeStorage) CreateReturns(result1 storage.File, result2 error) {
	fake.CreateStub = nil
	fake.createReturns = struct {
		result1 storage.File
		result2 error
	}{result1, result2}
}

func (fake *FakeStorage) Open(name string) (storage.File, error) {
	fake.openMutex.Lock()
	fake.openArgsForCall = append(fake.openArgsForCall, struct {
		name string
	}{name})
	fake.openMutex.Unlock()
	if fake.OpenStub != nil {
		return fake.OpenStub(name)
	} else {
		return fake.openReturns.result1, fake.openReturns.result2
	}
}

func (fake *FakeStorage) OpenCallCount() int {
	fake.openMutex.RLock()
	defer fake.openMutex.RUnlock()
	return len(fake.openArgsForCall)
}

func (fake *FakeStorage) OpenArgsForCall(i int) string {
	fake.openMutex.RLock()
	defer fake.openMutex.RUnlock()
	return fake.openArgsForCall[i].name
}

func (fake *FakeStorage) OpenReturns(result1 storage.File, result2 error) {
	fake.OpenStub = nil
	fake.openReturns = struct {
		result1 storage.File
		result2 error
	}{result1, result2}
}

func (fake *FakeStorage) WriteJson(name string, payload interface{}) error {
	fake.writeJsonMutex.Lock()
	fake.writeJsonArgsForCall = append(fake.writeJsonArgsForCall, struct {
		name    string
		payload interface{}
	}{name, payload})
	fake.writeJsonMutex.Unlock()
	if fake.WriteJsonStub != nil {
		return fake.WriteJsonStub(name, payload)
	} else {
		return fake.writeJsonReturns.result1
	}
}

func (fake *FakeStorage) WriteJsonCallCount() int {
	fake.writeJsonMutex.RLock()
	defer fake.writeJsonMutex.RUnlock()
	return len(fake.writeJsonArgsForCall)
}

func (fake *FakeStorage) WriteJsonArgsForCall(i int) (string, interface{}) {
	fake.writeJsonMutex.RLock()
	defer fake.writeJsonMutex.RUnlock()
	return fake.writeJsonArgsForCall[i].name, fake.writeJsonArgsForCall[i].payload
}

func (fake *FakeStorage) WriteJsonReturns(result1 error) {
	fake.WriteJsonStub = nil
	fake.writeJsonReturns = struct {
		result1 error
	}{result1}
}

var _ storage.Storage = new(FakeStorage)
