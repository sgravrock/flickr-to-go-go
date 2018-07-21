package storage_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/sgravrock/flickr-to-go-go/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"path"
)

var _ = Describe("Storage", func() {
	var dir string
	var subject Storage

	BeforeEach(func() {
		d, err := ioutil.TempDir("", "storage_test")
		dir = d
		Expect(err).To(BeNil())
		subject = NewFileStorage(d)
	})

	Describe("Create", func() {
		Context("When the parent directory doesn't exist", func() {
			It("should create the parent directory", func() {
				f, err := subject.Create("foo/bar")
				Expect(err).To(BeNil())
				f.Close()
				_, err = os.Stat(fmt.Sprintf("%s/foo/bar", dir))
				Expect(err).To(BeNil())
			})

			Context("When the creation fails", func() {
				BeforeEach(func() {
					// Create a regular file where we want a directory to be created
					f, err := subject.Create("foo")
					Expect(err).To(BeNil())
					f.Close()
				})

				It("should fail", func() {
					_, err := subject.Create("foo/bar")
					Expect(err).NotTo(BeNil())
				})
			})
		})
	})

	Describe("ListFiles", func() {
		It("returns the list of files in the specified directory", func() {
			err := os.Mkdir(path.Join(dir, "subdir"), 0777)
			Expect(err).To(BeNil())
			f, err := os.Create(path.Join(dir, "subdir/file1"))
			Expect(err).To(BeNil())
			f.Close()
			f, err = os.Create(path.Join(dir, "subdir/file2"))
			Expect(err).To(BeNil())
			f.Close()

			files, err := subject.ListFiles("subdir")
			Expect(err).To(BeNil())
			Expect(files).To(Equal([]string {"file1", "file2"}))
		})

		It("excludes directories", func() {
			err := os.MkdirAll(path.Join(dir, "subdir/subsub"), 0777)
			Expect(err).To(BeNil())
			files, err := subject.ListFiles("subdir")
			Expect(err).To(BeNil())
			Expect(len(files)).To(Equal(0))
		})
	})

	Describe("Move", func() {
		It("moves the file", func() {
			f, err := subject.Create("file")
			Expect(err).To(BeNil())
			f.Close()
			err = os.Mkdir(path.Join(dir, "subdir"), 0777)
			Expect(err).To(BeNil())

			err = subject.Move("file", "subdir/file")
			Expect(err).To(BeNil())

			_, err = os.Stat(path.Join(dir, "file"))
			Expect(err).NotTo(BeNil())
			_, err = os.Stat(path.Join(dir, "subdir/file"))
			Expect(err).To(BeNil())
		})
	})
})
