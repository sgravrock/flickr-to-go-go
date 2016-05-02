package storage

import (
	"io"
	"os"
	"path"
)

type Storage interface {
	Create(name string) (File, error)
	Open(name string) (File, error)
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

func (fs FileStorage) Create(name string) (File, error) {
	return os.Create(path.Join(fs.Rootdir, name))
}

func (fs FileStorage) Open(name string) (File, error) {
	return os.Open(path.Join(fs.Rootdir, name))
}
