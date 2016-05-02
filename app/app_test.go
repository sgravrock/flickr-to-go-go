package app_test

import (
	"bytes"
	"errors"
	"net/http"

	. "github.com/sgravrock/flickr-to-go-go/app"
	"github.com/sgravrock/flickr-to-go-go/auth/authfakes"
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
	)

	BeforeEach(func() {
		authenticator = new(authfakes.FakeAuthenticator)
		authenticator.AuthenticateReturns(new(http.Client), nil)
		fs = new(storagefakes.FakeStorage)
		ts = ghttp.NewServer()
		ts.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/services/rest/", "format=json&method=flickr.test.login&nojsoncallback=1"),
				ghttp.RespondWith(http.StatusNotFound, "nope"),
			),
		)
	})

	Context("When the destination directory doesn't exist", func() {
		BeforeEach(func() {
			fs.EnsureRootReturns(errors.New("stat foo: No such file or directory"))
		})

		It("prints an error and fails", func() {
			stderr := new(bytes.Buffer)
			ret := Run(ts.URL(), false, authenticator, fs, new(bytes.Buffer), stderr)
			Expect(ret).NotTo(Equal(0))
			Expect(authenticator.AuthenticateCallCount()).To(Equal(0))
			Expect(stderr.String()).To(ContainSubstring("foo: No such file or directory"))
		})
	})

	It("authenticates the user", func() {
		Run(ts.URL(), false, authenticator, fs, new(bytes.Buffer), new(bytes.Buffer))
		Expect(authenticator.AuthenticateCallCount()).To(Equal(1))
	})
})
