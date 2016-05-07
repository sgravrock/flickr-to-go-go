package storage_test

import (
	"fmt"
	"io/ioutil"
	"os"

	. "github.com/sgravrock/flickr-to-go-go/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
})
