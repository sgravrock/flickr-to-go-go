package storage

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"io/ioutil"
)

type Storage interface {
	EnsureRoot() error
	Exists(string) bool
	ListFiles(dir string) ([]string, error)
	Create(name string) (File, error)
	Open(name string) (File, error)
	Move(oldPath string, newPath string) error
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

func (fs FileStorage) Exists(name string) bool {
	_, err := os.Stat(path.Join(fs.Rootdir, name))
	return err == nil
}

func (fs FileStorage) ListFiles(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(path.Join(fs.Rootdir, dir))
	if err != nil {
		return nil, err
	}

	var result []string

	for _, file := range files {
		if !file.IsDir() {
			result = append(result, file.Name())
		}
	}

	return result, nil
}

func (fs FileStorage) Create(name string) (File, error) {
	p := path.Join(fs.Rootdir, name)
	err := ensureParentDir(p)
	if err != nil {
		return nil, err
	}
	return os.Create(p)
}

func (fs FileStorage) Open(name string) (File, error) {
	return os.Open(path.Join(fs.Rootdir, name))
}

func (fs FileStorage) Move(oldPath string, newPath string) error {
	return os.Rename(
		path.Join(fs.Rootdir, oldPath),
		path.Join(fs.Rootdir, newPath))
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

func ensureParentDir(p string) error {
	parent := path.Dir(p)
	return os.MkdirAll(parent, 0777)
}
