package app_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	. "github.com/sgravrock/flickr-to-go-go/app"
	"github.com/sgravrock/flickr-to-go-go/auth/authfakes"
	"github.com/sgravrock/flickr-to-go-go/clock/clockfakes"
	"github.com/sgravrock/flickr-to-go-go/dl/dlfakes"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage"

	"io"

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

	It("runs to completion even if the destination is empty", func() {
		Expect(stderr.String()).To(Equal(""))
		Expect(retval).To(Equal(0))
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
				})

				It("downloads all photo infos", func() {
					Expect(downloader.DownloadPhotoInfoCallCount()).To(Equal(2))
					_, _, id := downloader.DownloadPhotoInfoArgsForCall(0)
					Expect(id).To(Equal("123"))
					_, _, id = downloader.DownloadPhotoInfoArgsForCall(1)
					Expect(id).To(Equal("456"))
				})

				Context("And all photo infos are present", func() {
					BeforeEach(func() {
						downloader.PhotoInfoExistsReturns(true)
					})

					It("downloads info for recently created or updated photos", func() {
						Expect(downloader.DownloadPhotoInfoCallCount()).To(Equal(1))
						_, _, id := downloader.DownloadPhotoInfoArgsForCall(0)
						Expect(id).To(Equal("456"))
					})
				})

				Context("And some photo infos are missing", func() {
					BeforeEach(func() {
						downloader.PhotoInfoExistsReturns(false)
					})

					It("downloads missing infos", func() {
						Expect(downloader.DownloadPhotoInfoCallCount()).To(Equal(2))
						_, _, id := downloader.DownloadPhotoInfoArgsForCall(0)
						Expect(id).To(Equal("123"))
						_, _, id = downloader.DownloadPhotoInfoArgsForCall(1)
						Expect(id).To(Equal("456"))
					})
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

		Context("When a previously downloaded photo has been deleted", func() {
			var (
				oldInfoPath string
				oldOrigPath string
			)

			BeforeEach(func() {
				infoDir := path.Join(dir, "photo-info")
				err := os.MkdirAll(infoDir, 0777)
				Expect(err).To(BeNil())
				oldInfoPath = path.Join(infoDir, "678.json")
				infoFile, err := os.Create(oldInfoPath)
				Expect(err).To(BeNil())
				defer infoFile.Close()
				io.WriteString(infoFile, "info")

				_, err = os.Stat(oldInfoPath)
				Expect(err).To(BeNil())

				origDir := path.Join(dir, "originals")
				err = os.MkdirAll(origDir, 0777)
				Expect(err).To(BeNil())
				oldOrigPath = path.Join(origDir, "678.jpg")
				origFile, err := os.Create(oldOrigPath)
				Expect(err).To(BeNil())
				defer origFile.Close()
				io.WriteString(origFile, "orig")
			})

			It("moves the info file to the attic", func() {
				_, err := os.Stat(oldInfoPath)
				Expect(err).NotTo(BeNil())

				newInfoPath := path.Join(dir, "attic/photo-info/678.json")
				b, err := ioutil.ReadFile(newInfoPath)
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal("info"))
				Expect(stdout.String()).To(ContainSubstring(
					"Moving info of deleted photo 678 to attic"))
			})

			It("moves the original to the attic", func() {
				_, err := os.Stat(oldOrigPath)
				Expect(err).NotTo(BeNil())

				newOrigPath := path.Join(dir, "attic/originals/678.jpg")
				b, err := ioutil.ReadFile(newOrigPath)
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal("orig"))
				Expect(stdout.String()).To(ContainSubstring(
					"Moving original of deleted photo 678 to attic"))

			})
		})

		Context("When a previously downloaded photo has not been updated", func() {
			var (
				infoPath string
				origPath string
			)

			BeforeEach(func() {
				infoDir := path.Join(dir, "photo-info")
				err := os.MkdirAll(infoDir, 0777)
				Expect(err).To(BeNil())
				infoPath = path.Join(infoDir, "123.json")
				infoFile, err := os.Create(infoPath)
				Expect(err).To(BeNil())
				defer infoFile.Close()
				io.WriteString(infoFile, "info")

				_, err = os.Stat(infoPath)
				Expect(err).To(BeNil())

				origDir := path.Join(dir, "originals")
				err = os.MkdirAll(origDir, 0777)
				Expect(err).To(BeNil())
				origPath = path.Join(origDir, "123.jpg")
				origFile, err := os.Create(origPath)
				Expect(err).To(BeNil())
				defer origFile.Close()
				io.WriteString(origFile, "orig")
			})

			It("does not move the info file to the attic", func() {
				b, err := ioutil.ReadFile(infoPath)
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal("info"))
			})

			It("does not move the original to the attic", func() {
				b, err := ioutil.ReadFile(origPath)
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal("orig"))
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

		It("exits successfully", func() {
			Expect(stderr.String()).To(Equal(""))
			Expect(retval).To(Equal(0))
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
