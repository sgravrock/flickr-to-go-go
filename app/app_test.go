package app_test

import (
	"bytes"
	"errors"
	"net/http"

	. "github.com/sgravrock/flickr-to-go-go/app"
	"github.com/sgravrock/flickr-to-go-go/auth/authfakes"
	"github.com/sgravrock/flickr-to-go-go/dl/dlfakes"
	"github.com/sgravrock/flickr-to-go-go/flickrapi"
	"github.com/sgravrock/flickr-to-go-go/storage/storagefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("App", func() {
	var (
		authenticator *authfakes.FakeAuthenticator
		ts            *ghttp.Server
		fs            *storagefakes.FakeStorage
		downloader    *dlfakes.FakeDownloader
		stdout        *bytes.Buffer
		stderr        *bytes.Buffer
		retval        int
	)

	BeforeEach(func() {
		authenticator = new(authfakes.FakeAuthenticator)
		authenticator.AuthenticateReturns(new(http.Client), nil)
		fs = new(storagefakes.FakeStorage)
		downloader = new(dlfakes.FakeDownloader)
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
		retval = Run(ts.URL(), false, authenticator, downloader, fs, stdout, stderr)
	})

	Context("When the destination directory doesn't exist", func() {
		BeforeEach(func() {
			fs.EnsureRootReturns(errors.New("stat foo: No such file or directory"))
		})

		It("prints an error and fails", func() {
			Expect(retval).NotTo(Equal(0))
			Expect(authenticator.AuthenticateCallCount()).To(Equal(0))
			Expect(stderr.String()).To(ContainSubstring("foo: No such file or directory"))
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

		It("downloads the photos' info", func() {
			Expect(downloader.DownloadPhotoInfoCallCount()).To(Equal(2))
			_, _, id := downloader.DownloadPhotoInfoArgsForCall(0)
			Expect(id).To(Equal("123"))
			_, _, id = downloader.DownloadPhotoInfoArgsForCall(1)
			Expect(id).To(Equal("456"))
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

		It("downloads the original photos", func() {
			Expect(downloader.DownloadOriginalCallCount()).To(Equal(2))
			_, _, p1 := downloader.DownloadOriginalArgsForCall(0)
			Expect(p1.Id()).To(Equal("123"))
			_, _, p2 := downloader.DownloadOriginalArgsForCall(1)
			Expect(p2.Id()).To(Equal("456"))
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
})
