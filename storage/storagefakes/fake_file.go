// This file was generated by counterfeiter
package storagefakes

import (
	"os"
	"sync"

	"github.com/sgravrock/flickr-to-go-go/storage"
)

type FakeFile struct {
	CloseStub        func() error
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	closeReturns     struct {
		result1 error
	}
	ReadStub        func(p []byte) (n int, err error)
	readMutex       sync.RWMutex
	readArgsForCall []struct {
		p []byte
	}
	readReturns struct {
		result1 int
		result2 error
	}
	ReadAtStub        func(p []byte, off int64) (n int, err error)
	readAtMutex       sync.RWMutex
	readAtArgsForCall []struct {
		p   []byte
		off int64
	}
	readAtReturns struct {
		result1 int
		result2 error
	}
	WriteStub        func(p []byte) (n int, err error)
	writeMutex       sync.RWMutex
	writeArgsForCall []struct {
		p []byte
	}
	writeReturns struct {
		result1 int
		result2 error
	}
	WriteAtStub        func(p []byte, off int64) (n int, err error)
	writeAtMutex       sync.RWMutex
	writeAtArgsForCall []struct {
		p   []byte
		off int64
	}
	writeAtReturns struct {
		result1 int
		result2 error
	}
	SeekStub        func(offset int64, whence int) (int64, error)
	seekMutex       sync.RWMutex
	seekArgsForCall []struct {
		offset int64
		whence int
	}
	seekReturns struct {
		result1 int64
		result2 error
	}
	StatStub        func() (os.FileInfo, error)
	statMutex       sync.RWMutex
	statArgsForCall []struct{}
	statReturns     struct {
		result1 os.FileInfo
		result2 error
	}
}

func (fake *FakeFile) Close() error {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		return fake.CloseStub()
	} else {
		return fake.closeReturns.result1
	}
}

func (fake *FakeFile) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeFile) CloseReturns(result1 error) {
	fake.CloseStub = nil
	fake.closeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeFile) Read(p []byte) (n int, err error) {
	var pCopy []byte
	if p != nil {
		pCopy = make([]byte, len(p))
		copy(pCopy, p)
	}
	fake.readMutex.Lock()
	fake.readArgsForCall = append(fake.readArgsForCall, struct {
		p []byte
	}{pCopy})
	fake.readMutex.Unlock()
	if fake.ReadStub != nil {
		return fake.ReadStub(p)
	} else {
		return fake.readReturns.result1, fake.readReturns.result2
	}
}

func (fake *FakeFile) ReadCallCount() int {
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	return len(fake.readArgsForCall)
}

func (fake *FakeFile) ReadArgsForCall(i int) []byte {
	fake.readMutex.RLock()
	defer fake.readMutex.RUnlock()
	return fake.readArgsForCall[i].p
}

func (fake *FakeFile) ReadReturns(result1 int, result2 error) {
	fake.ReadStub = nil
	fake.readReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeFile) ReadAt(p []byte, off int64) (n int, err error) {
	var pCopy []byte
	if p != nil {
		pCopy = make([]byte, len(p))
		copy(pCopy, p)
	}
	fake.readAtMutex.Lock()
	fake.readAtArgsForCall = append(fake.readAtArgsForCall, struct {
		p   []byte
		off int64
	}{pCopy, off})
	fake.readAtMutex.Unlock()
	if fake.ReadAtStub != nil {
		return fake.ReadAtStub(p, off)
	} else {
		return fake.readAtReturns.result1, fake.readAtReturns.result2
	}
}

func (fake *FakeFile) ReadAtCallCount() int {
	fake.readAtMutex.RLock()
	defer fake.readAtMutex.RUnlock()
	return len(fake.readAtArgsForCall)
}

func (fake *FakeFile) ReadAtArgsForCall(i int) ([]byte, int64) {
	fake.readAtMutex.RLock()
	defer fake.readAtMutex.RUnlock()
	return fake.readAtArgsForCall[i].p, fake.readAtArgsForCall[i].off
}

func (fake *FakeFile) ReadAtReturns(result1 int, result2 error) {
	fake.ReadAtStub = nil
	fake.readAtReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeFile) Write(p []byte) (n int, err error) {
	var pCopy []byte
	if p != nil {
		pCopy = make([]byte, len(p))
		copy(pCopy, p)
	}
	fake.writeMutex.Lock()
	fake.writeArgsForCall = append(fake.writeArgsForCall, struct {
		p []byte
	}{pCopy})
	fake.writeMutex.Unlock()
	if fake.WriteStub != nil {
		return fake.WriteStub(p)
	} else {
		return fake.writeReturns.result1, fake.writeReturns.result2
	}
}

func (fake *FakeFile) WriteCallCount() int {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return len(fake.writeArgsForCall)
}

func (fake *FakeFile) WriteArgsForCall(i int) []byte {
	fake.writeMutex.RLock()
	defer fake.writeMutex.RUnlock()
	return fake.writeArgsForCall[i].p
}

func (fake *FakeFile) WriteReturns(result1 int, result2 error) {
	fake.WriteStub = nil
	fake.writeReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeFile) WriteAt(p []byte, off int64) (n int, err error) {
	var pCopy []byte
	if p != nil {
		pCopy = make([]byte, len(p))
		copy(pCopy, p)
	}
	fake.writeAtMutex.Lock()
	fake.writeAtArgsForCall = append(fake.writeAtArgsForCall, struct {
		p   []byte
		off int64
	}{pCopy, off})
	fake.writeAtMutex.Unlock()
	if fake.WriteAtStub != nil {
		return fake.WriteAtStub(p, off)
	} else {
		return fake.writeAtReturns.result1, fake.writeAtReturns.result2
	}
}

func (fake *FakeFile) WriteAtCallCount() int {
	fake.writeAtMutex.RLock()
	defer fake.writeAtMutex.RUnlock()
	return len(fake.writeAtArgsForCall)
}

func (fake *FakeFile) WriteAtArgsForCall(i int) ([]byte, int64) {
	fake.writeAtMutex.RLock()
	defer fake.writeAtMutex.RUnlock()
	return fake.writeAtArgsForCall[i].p, fake.writeAtArgsForCall[i].off
}

func (fake *FakeFile) WriteAtReturns(result1 int, result2 error) {
	fake.WriteAtStub = nil
	fake.writeAtReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeFile) Seek(offset int64, whence int) (int64, error) {
	fake.seekMutex.Lock()
	fake.seekArgsForCall = append(fake.seekArgsForCall, struct {
		offset int64
		whence int
	}{offset, whence})
	fake.seekMutex.Unlock()
	if fake.SeekStub != nil {
		return fake.SeekStub(offset, whence)
	} else {
		return fake.seekReturns.result1, fake.seekReturns.result2
	}
}

func (fake *FakeFile) SeekCallCount() int {
	fake.seekMutex.RLock()
	defer fake.seekMutex.RUnlock()
	return len(fake.seekArgsForCall)
}

func (fake *FakeFile) SeekArgsForCall(i int) (int64, int) {
	fake.seekMutex.RLock()
	defer fake.seekMutex.RUnlock()
	return fake.seekArgsForCall[i].offset, fake.seekArgsForCall[i].whence
}

func (fake *FakeFile) SeekReturns(result1 int64, result2 error) {
	fake.SeekStub = nil
	fake.seekReturns = struct {
		result1 int64
		result2 error
	}{result1, result2}
}

func (fake *FakeFile) Stat() (os.FileInfo, error) {
	fake.statMutex.Lock()
	fake.statArgsForCall = append(fake.statArgsForCall, struct{}{})
	fake.statMutex.Unlock()
	if fake.StatStub != nil {
		return fake.StatStub()
	} else {
		return fake.statReturns.result1, fake.statReturns.result2
	}
}

func (fake *FakeFile) StatCallCount() int {
	fake.statMutex.RLock()
	defer fake.statMutex.RUnlock()
	return len(fake.statArgsForCall)
}

func (fake *FakeFile) StatReturns(result1 os.FileInfo, result2 error) {
	fake.StatStub = nil
	fake.statReturns = struct {
		result1 os.FileInfo
		result2 error
	}{result1, result2}
}

var _ storage.File = new(FakeFile)
