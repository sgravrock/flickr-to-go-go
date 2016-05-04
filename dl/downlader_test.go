package dl_test

import (
	"errors"

	. "github.com/sgravrock/flickr-to-go-go/dl"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/flickrapi/flickrapifakes"
	"github.com/sgravrock/flickr-to-go-go/storage/storagefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Downloader", func() {
	var subject Downloader
	var flickrClient *flickrapifakes.FakeClient
	var fs *storagefakes.FakeStorage

	BeforeEach(func() {
		flickrClient = new(flickrapifakes.FakeClient)
		fs = new(storagefakes.FakeStorage)
		subject = NewDownloader()
	})

	Describe("DownloadPhotolist", func() {
		var err error
		var result []flickrapi.PhotoListEntry

		JustBeforeEach(func() {
			result, err = subject.DownloadPhotolist(flickrClient, fs)
		})

		It("should request the photo list", func() {
			Expect(flickrClient.GetPhotosCallCount()).To(Equal(1))
		})

		Context("When the photo list fetch fails", func() {
			BeforeEach(func() {
				flickrClient.GetPhotosReturns(nil, errors.New("nope"))
			})

			It("should fail", func() {
				Expect(result).To(BeNil())
				Expect(err).NotTo(BeNil())
			})

			It("should not save anything", func() {
				Expect(fs.WriteJsonCallCount()).To(Equal(0))
			})
		})

		Context("When the photo list fetch succeeds", func() {
			var expected []flickrapi.PhotoListEntry
			BeforeEach(func() {
				expected = []flickrapi.PhotoListEntry{
					flickrapi.PhotoListEntry{
						Id:       "123",
						Owner:    "1234@N02",
						Secret:   "asdf",
						Server:   "1518",
						Farm:     2,
						Title:    "t1",
						Ispublic: 1,
						Isfriend: 0,
						Isfamily: 0,
					},
				}
				flickrClient.GetPhotosReturns(expected, nil)
			})

			It("should save the photo list", func() {
				Expect(fs.WriteJsonCallCount()).To(Equal(1))
				name, payload := fs.WriteJsonArgsForCall(0)
				Expect(name).To(Equal("photolist.json"))
				Expect(payload).To(Equal(expected))
			})

			Context("When the save fails", func() {
				BeforeEach(func() {
					fs.WriteJsonReturns(errors.New("nope!"))
				})

				It("should fail", func() {
					Expect(result).To(BeNil())
					Expect(err).NotTo(BeNil())
				})
			})

			Context("When the save succeeds", func() {
				It("should return the photos", func() {
					Expect(result).To(Equal(expected))
					Expect(err).To(BeNil())
				})
			})
		})
	})

	Describe("DownloadPhotoInfo", func() {
		var err error

		JustBeforeEach(func() {
			err = subject.DownloadPhotoInfo(flickrClient, fs, "789")
		})

		It("requests the photo info", func() {
			Expect(flickrClient.GetPhotoInfoCallCount()).To(Equal(1))
			Expect(flickrClient.GetPhotoInfoArgsForCall(0)).To(Equal("789"))
		})

		Context("When the request succeeds", func() {
			var photo flickrapi.PhotoInfo

			BeforeEach(func() {
				photo = flickrapi.PhotoInfo{Id: "789"}
				flickrClient.GetPhotoInfoReturns(photo, nil)
			})

			It("saves the photo", func() {
				Expect(fs.WriteJsonCallCount()).To(Equal(1))
				name, payload := fs.WriteJsonArgsForCall(0)
				Expect(name).To(Equal("photo-info/789.json"))
				Expect(payload).To(Equal(photo))
			})
		})
	})
})
