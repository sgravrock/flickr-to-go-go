package storage

import (
	"encoding/json"
	"io"
	"os"
	"path"
)

type Storage interface {
	EnsureRoot() error
	Create(name string) (File, error)
	Open(name string) (File, error)
	WriteJson(name string, payload interface{}) error
}

type File interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
	io.Seeker
	Stat() (os.FileInfo, error)
}

type FileStorage struct{ Rootdir string }

func NewFileStorage(rootDir string) Storage {
	return &FileStorage{rootDir}
}

func (fs FileStorage) EnsureRoot() error {
	_, err := os.Stat(fs.Rootdir)
	return err
}

func (fs FileStorage) Create(name string) (File, error) {
	return os.Create(path.Join(fs.Rootdir, name))
}

func (fs FileStorage) Open(name string) (File, error) {
	return os.Open(path.Join(fs.Rootdir, name))
}

func (fs FileStorage) WriteJson(name string, payload interface{}) error {
	data, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}
	f, err := fs.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}
