package dl_test

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/sgravrock/flickr-to-go-go/dl"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/flickrapi/flickrapifakes"
	"github.com/sgravrock/flickr-to-go-go/storage/storagefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
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
			var photolist []flickrapi.PhotoListEntry
			var rawPhotos []map[string]interface{}
			BeforeEach(func() {
				rawPhotos = []map[string]interface{}{
					map[string]interface{}{"id": "123"},
				}
				photolist = []flickrapi.PhotoListEntry{
					flickrapi.PhotoListEntry{
						Data: rawPhotos[0],
					},
				}
				flickrClient.GetPhotosReturns(photolist, nil)
			})

			It("should save the photo list", func() {
				Expect(fs.WriteJsonCallCount()).To(Equal(1))
				name, payload := fs.WriteJsonArgsForCall(0)
				Expect(name).To(Equal("photolist.json"))
				Expect(payload).To(Equal(rawPhotos))
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
					Expect(result).To(Equal(photolist))
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

		Context("When the photo info already exists", func() {
			BeforeEach(func() {
				fs.ExistsStub = func(name string) bool {
					return name == "photo-info/789.json"
				}
			})

			It("does not download anything", func() {
				Expect(flickrClient.GetPhotoInfoCallCount()).To(Equal(0))
				Expect(fs.WriteJsonCallCount()).To(Equal(0))
			})
		})

		It("requests the photo info", func() {
			Expect(flickrClient.GetPhotoInfoCallCount()).To(Equal(1))
			Expect(flickrClient.GetPhotoInfoArgsForCall(0)).To(Equal("789"))
		})

		Context("When the request fails", func() {
			BeforeEach(func() {
				flickrClient.GetPhotoInfoReturns(nil, errors.New("nope"))
			})

			It("fails", func() {
				Expect(err).NotTo(BeNil())
				Expect(fs.WriteJsonCallCount()).To(Equal(0))
			})
		})

		Context("When the request succeeds", func() {
			var photo map[string]interface{}

			BeforeEach(func() {
				photo = map[string]interface{}{"id": "789"}
				flickrClient.GetPhotoInfoReturns(photo, nil)
			})

			It("saves the photo info", func() {
				Expect(fs.WriteJsonCallCount()).To(Equal(1))
				name, payload := fs.WriteJsonArgsForCall(0)
				Expect(name).To(Equal("photo-info/789.json"))
				Expect(payload).To(Equal(photo))
			})

			Context("When the save fails", func() {
				BeforeEach(func() {
					fs.WriteJsonReturns(errors.New("nope"))
				})

				It("fails", func() {
					Expect(err).NotTo(BeNil())
				})
			})
		})
	})

	Describe("DownloadOriginal", func() {
		var photo flickrapi.PhotoListEntry
		var ts *ghttp.Server
		var path string
		var err error

		BeforeEach(func() {
			ts = ghttp.NewServer()
			path = "/path/some-number_o.jpg"
			photo = flickrapi.PhotoListEntry{
				Data: map[string]interface{}{
					"id":    "12345",
					"url_o": fmt.Sprintf("%s%s", ts.URL(), path),
				},
			}
		})

		JustBeforeEach(func() {
			err = subject.DownloadOriginal(new(http.Client), fs, photo)
		})

		Context("When the photo already exists", func() {
			BeforeEach(func() {
				fs.ExistsStub = func(name string) bool {
					return name == "originals/12345.jpg"
				}
			})

			It("does not download anything", func() {
				Expect(fs.CreateCallCount()).To(Equal(0))
			})
		})

		Context("When the request returns a non-200 response", func() {
			BeforeEach(func() {
				ts.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", path, ""),
						ghttp.RespondWith(403, "Not yours"),
					),
				)
			})

			It("fails", func() {
				Expect(err).NotTo(BeNil())
			})
		})

		Context("When the request succeeds", func() {
			var contents []byte
			var file *storagefakes.FakeFile

			BeforeEach(func() {
				contents = []byte("\xff\xd8\xff\xe1\x16&Exif\x00\x00II*\x00\x08\x00\x00\x00")
				ts.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", path, ""),
						ghttp.RespondWith(200, contents),
					),
				)
				file = new(storagefakes.FakeFile)
				fs.CreateReturns(file, nil)
			})

			It("saves the photo", func() {
				Expect(fs.CreateCallCount()).To(Equal(1))
				Expect(fs.CreateArgsForCall(0)).To(Equal("originals/12345.jpg"))
				Expect(file.WriteCallCount()).To(Equal(1))
				Expect(file.WriteArgsForCall(0)).To(Equal(contents))
				Expect(file.CloseCallCount()).To(Equal(1))
			})

			Context("When the open fails", func() {
				BeforeEach(func() {
					fs.CreateReturns(nil, errors.New("nope"))
				})

				It("fails", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("When the write fails", func() {
				BeforeEach(func() {
					file.WriteReturns(-1, errors.New("nope"))
				})

				It("fails", func() {
					Expect(err).NotTo(BeNil())
				})
			})
		})
	})
})
