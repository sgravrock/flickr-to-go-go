package app_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	. "github.com/sgravrock/flickr-to-go-go/app"
	"github.com/sgravrock/flickr-to-go-go/auth/authfakes"
	"github.com/sgravrock/flickr-to-go-go/clock/clockfakes"
	"github.com/sgravrock/flickr-to-go-go/dl/dlfakes"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("App", func() {
	var (
		authenticator *authfakes.FakeAuthenticator
		ts            *ghttp.Server
		fs            storage.Storage
		dir           string
		downloader    *dlfakes.FakeDownloader
		clock         *clockfakes.FakeClock
		stdout        *bytes.Buffer
		stderr        *bytes.Buffer
		retval        int
	)

	BeforeEach(func() {
		authenticator = new(authfakes.FakeAuthenticator)
		authenticator.AuthenticateReturns(new(http.Client), nil)
		d, err := ioutil.TempDir("", "storage_test")
		dir = d
		Expect(err).To(BeNil())
		fs = storage.NewFileStorage(dir)
		downloader = new(dlfakes.FakeDownloader)
		clock = new(clockfakes.FakeClock)
		stdout = new(bytes.Buffer)
		stderr = new(bytes.Buffer)
		ts = ghttp.NewServer()
		ts.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/services/rest/", "format=json&method=flickr.test.login&nojsoncallback=1"),
				ghttp.RespondWith(http.StatusNotFound, "nope"),
			),
		)
	})

	JustBeforeEach(func() {
		retval = Run(ts.URL(), false, authenticator, downloader, fs, clock,
			stdout, stderr)
	})

	Context("When the destination directory doesn't exist", func() {
		BeforeEach(func() {
			err := os.Remove(dir)
			Expect(err).To(BeNil())
		})

		It("prints an error and fails", func() {
			Expect(retval).NotTo(Equal(0))
			Expect(authenticator.AuthenticateCallCount()).To(Equal(0))
			Expect(stderr.String()).To(ContainSubstring("no such file or directory"))
		})
	})

	It("authenticates the user", func() {
		Expect(authenticator.AuthenticateCallCount()).To(Equal(1))
	})

	It("downloads the photo list", func() {
		Expect(downloader.DownloadPhotolistCallCount()).To(Equal(1))
	})

	Context("When the photo list download fails", func() {
		BeforeEach(func() {
			downloader.DownloadPhotolistReturns(nil, errors.New("nope"))
		})

		It("prints an error and fails", func() {
			Expect(retval).NotTo(Equal(0))
			Expect(stderr.String()).To(ContainSubstring("Error downloading photo list: nope"))
		})
	})

	Context("When the photo list download succeeds", func() {
		var photos []flickrapi.PhotoListEntry

		BeforeEach(func() {
			photos = []flickrapi.PhotoListEntry{
				flickrapi.PhotoListEntry{Data: map[string]interface{}{"id": "123"}},
				flickrapi.PhotoListEntry{Data: map[string]interface{}{"id": "456"}},
			}
			downloader.DownloadPhotolistReturns(photos, nil)
		})

		Context("When the timestamp file doesn't exist", func() {
			It("downloads all photos in the list", func() {
				Expect(downloader.DownloadPhotoInfoCallCount()).To(Equal(2))
				_, _, id := downloader.DownloadPhotoInfoArgsForCall(0)
				Expect(id).To(Equal("123"))
				_, _, id = downloader.DownloadPhotoInfoArgsForCall(1)
				Expect(id).To(Equal("456"))

				Expect(downloader.DownloadOriginalCallCount()).To(Equal(2))
				_, _, p1 := downloader.DownloadOriginalArgsForCall(0)
				Expect(p1.Id()).To(Equal("123"))
				_, _, p2 := downloader.DownloadOriginalArgsForCall(1)
				Expect(p2.Id()).To(Equal("456"))
			})
		})

		Context("When the timestamp file exists", func() {
			BeforeEach(func() {
				err := ioutil.WriteFile(dir+"/timestamp", []byte("1257894000\n"), 0600)
				Expect(err).To(BeNil())
			})

			It("gets the list of photos created or modified since the timestamp", func() {
				Expect(downloader.GetRecentPhotoIdsCallCount()).To(Equal(1))
				timestamp, _ := downloader.GetRecentPhotoIdsArgsForCall(0)
				Expect(timestamp).To(Equal(uint32(1257894000)))
			})

			Context("When the GetRecentPhotoIds call fails", func() {
				BeforeEach(func() {
					downloader.GetRecentPhotoIdsReturns(nil, errors.New("nope"))
				})

				It("downloads all photos in the list", func() {
					Expect(downloader.DownloadPhotoInfoCallCount()).To(Equal(2))
					_, _, id := downloader.DownloadPhotoInfoArgsForCall(0)
					Expect(id).To(Equal("123"))
					_, _, id = downloader.DownloadPhotoInfoArgsForCall(1)
					Expect(id).To(Equal("456"))

					Expect(downloader.DownloadOriginalCallCount()).To(Equal(2))
					_, _, p1 := downloader.DownloadOriginalArgsForCall(0)
					Expect(p1.Id()).To(Equal("123"))
					_, _, p2 := downloader.DownloadOriginalArgsForCall(1)
					Expect(p2.Id()).To(Equal("456"))
				})
			})

			Context("When the GetRecentPhotoIds call succeeds", func() {
				BeforeEach(func() {
					downloader.GetRecentPhotoIdsReturns([]string{"456"}, nil)
				});

				It("downloads all photo infos", func() {
					Expect(downloader.DownloadPhotoInfoCallCount()).To(Equal(2))
					_, _, id := downloader.DownloadPhotoInfoArgsForCall(0)
					Expect(id).To(Equal("123"))
					_, _, id = downloader.DownloadPhotoInfoArgsForCall(1)
					Expect(id).To(Equal("456"))
				})

				Context("And all originals are present", func() {
					BeforeEach(func() {
						downloader.OriginalExistsReturns(true)
					})

					It("downloads originals of recently created or updated photos", func() {
						Expect(downloader.DownloadOriginalCallCount()).To(Equal(1))
						_, _, p2 := downloader.DownloadOriginalArgsForCall(0)
						Expect(p2.Id()).To(Equal("456"))
					})

				})

				Context("And some originals are missing", func() {
					BeforeEach(func() {
						downloader.OriginalExistsReturns(false)
					})

					It("downloads other originals that are missing", func() {
						Expect(downloader.DownloadOriginalCallCount()).To(Equal(2))
						_, _, p1 := downloader.DownloadOriginalArgsForCall(0)
						Expect(p1.Id()).To(Equal("123"))
						_, _, p2 := downloader.DownloadOriginalArgsForCall(1)
						Expect(p2.Id()).To(Equal("456"))
					})
				})
			})
		})

		Context("When a photo info download fails", func() {
			BeforeEach(func() {
				downloader.DownloadPhotoInfoReturns(errors.New("nope"))
			})

			It("prints an error and fails", func() {
				Expect(retval).NotTo(Equal(0))
				Expect(stderr.String()).To(ContainSubstring(
					"Error downloading info for 123: nope"))
			})
		})

		Context("When an original fails to download", func() {
			BeforeEach(func() {
				downloader.DownloadOriginalReturns(errors.New("nope"))
			})

			It("prints an error and fails", func() {
				Expect(retval).NotTo(Equal(0))
				Expect(stderr.String()).To(ContainSubstring(
					"Error downloading original for 123: nope"))
			})
		})
	})

	Context("When everything succeeds", func() {
		BeforeEach(func() {
			photos := []flickrapi.PhotoListEntry{
				flickrapi.PhotoListEntry{Data: map[string]interface{}{"id": "456"}},
			}
			downloader.DownloadPhotolistReturns(photos, nil)
			clock.NowReturns(time.Unix(1257894000, 0))
		})

		It("writes a timestamp", func() {
			contents, err := ioutil.ReadFile(dir + "/timestamp")
			Expect(err).To(BeNil())
			Expect(contents).To(Equal([]byte("1257894000\n")))
		})
	})

	Context("When something fails", func() {
		BeforeEach(func() {
			photos := []flickrapi.PhotoListEntry{
				flickrapi.PhotoListEntry{Data: map[string]interface{}{"id": "456"}},
			}
			downloader.DownloadPhotolistReturns(photos, nil)
			downloader.DownloadOriginalReturns(errors.New("nope"))
		})

		It("doesn't write a timestamp", func() {
			_, err := os.Stat(dir + "/timestamp")
			_, ok := err.(*os.PathError)
			Expect(ok).To(BeTrue())
		})
	})
})
