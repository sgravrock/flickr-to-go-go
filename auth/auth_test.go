package auth_test

import (
	"errors"
	"net/http"
	"os"

	. "github.com/sgravrock/flickr-to-go-go/auth"

	"io/ioutil"

	"encoding/gob"

	"github.com/mrjones/oauth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sgravrock/flickr-to-go-go/auth/authfakes"
	"github.com/sgravrock/flickr-to-go-go/storage"
	"github.com/sgravrock/flickr-to-go-go/storage/storagefakes"
)

var _ = Describe("Authenticate", func() {
	var subject Authenticator
	var oauthClient authfakes.FakeOauthClient
	var oauthConsumer *authfakes.FakeOauthConsumer
	var ui *authfakes.FakeUiAdapter
	var fs *storagefakes.FakeStorage
	var credsFile *storagefakes.FakeFile
	var credsCreateError error

	BeforeEach(func() {
		oauthClient = *new(authfakes.FakeOauthClient)
		oauthConsumer = new(authfakes.FakeOauthConsumer)
		ui = new(authfakes.FakeUiAdapter)
		oauthClient.NewConsumerReturns(oauthConsumer)
		fs = new(storagefakes.FakeStorage)
		credsFile = nil
		credsCreateError = errors.New("not mocked")
		fs.CreateStub = func(name string) (storage.File, error) {
			if name == "flickr-credentials" {
				return credsFile, credsCreateError
			}

			return nil, errors.New("not mocked")
		}
		fs.OpenStub = func(name string) (storage.File, error) {
			return nil, errors.New("not mocked")
		}

		subject = NewAuthenticator("token", "secret", fs, &oauthClient, ui)
	})

	Context("When there are saved credentials", func() {
		expectedToken := oauth.AccessToken{"token", "secret", nil}
		var httpClient http.Client
		var result *http.Client
		var tempPath string

		BeforeEach(func() {
			tempFile, err := ioutil.TempFile("", "creds")
			Expect(err).To(BeNil())
			defer tempFile.Close()
			tempPath = tempFile.Name()
			gob.NewEncoder(tempFile).Encode(expectedToken)
			tempFile.Seek(0, 0)
			fs.OpenStub = func(name string) (storage.File, error) {
				if name == "flickr-credentials" {
					return os.Open(tempPath)
				} else {
					return nil, errors.New("not mocked")
				}
			}
			httpClient = http.Client{}
			oauthConsumer.MakeHttpClientReturns(&httpClient, nil)
			result, _ = subject.Authenticate(false)
		})

		AfterEach(func() {
			os.Remove(tempPath)
		})

		It("should return the created HTTP client", func() {
			Expect(oauthConsumer.MakeHttpClientCallCount()).To(Equal(1))
			Expect(result).To(BeIdenticalTo(&httpClient))
		})

		It("should not prompt the user to get a code", func() {
			Expect(ui.PromptForAccessCodeCallCount()).To(Equal(0))
		})
	})

	Context("When obtaining a request token fails", func() {
		var tokenError error

		BeforeEach(func() {
			tokenError = errors.New("nope")
			oauthConsumer.GetRequestTokenAndUrlReturns(nil, "", tokenError)
		})

		It("should return the error", func() {
			t, err := subject.Authenticate(true)
			Expect(t).To(BeNil())
			Expect(err).To(BeIdenticalTo(tokenError))
		})
	})

	Context("When a request token is obtained", func() {
		var requestToken oauth.RequestToken

		BeforeEach(func() {
			requestToken = oauth.RequestToken{"the request token", ""}
			oauthConsumer.GetRequestTokenAndUrlReturns(&requestToken,
				"the://url", nil)
		})

		It("should prompt the user for authentication", func() {
			subject.Authenticate(false)
			Expect(ui.PromptForAccessCodeCallCount()).To(Equal(1))
			Expect(ui.PromptForAccessCodeArgsForCall(0)).To(Equal("the://url"))
		})

		Context("When obtaining an access code fails", func() {
			var codeError error

			BeforeEach(func() {
				codeError = errors.New("nope")
				ui.PromptForAccessCodeReturns("", codeError)
			})

			It("should return the error", func() {
				t, err := subject.Authenticate(true)
				Expect(t).To(BeNil())
				Expect(err).To(BeIdenticalTo(codeError))
			})
		})

		Context("When the user enters an access code", func() {
			BeforeEach(func() {
				ui.PromptForAccessCodeReturns("the code", nil)
			})

			It("should authorize the request token and code", func() {
				subject.Authenticate(true)
				Expect(oauthConsumer.AuthorizeTokenCallCount()).To(Equal(1))
				token, code := oauthConsumer.AuthorizeTokenArgsForCall(0)
				Expect(token).To(Equal(&requestToken))
				Expect(code).To(Equal("the code"))
			})

			Context("When authorization fails", func() {
				var authError error

				BeforeEach(func() {
					authError = errors.New("nope")
					oauthConsumer.AuthorizeTokenReturns(nil, authError)
				})

				It("should return the error", func() {
					t, err := subject.Authenticate(true)
					Expect(t).To(BeNil())
					Expect(err).To(BeIdenticalTo(authError))
				})

			})

			Context("When authorization succeeds", func() {
				var accessToken oauth.AccessToken
				var httpClient http.Client

				BeforeEach(func() {
					accessToken = oauth.AccessToken{
						"access token", "", nil,
					}
					oauthConsumer.AuthorizeTokenReturns(&accessToken, nil)
					httpClient = http.Client{}
					oauthConsumer.MakeHttpClientReturns(&httpClient, nil)
				})

				It("should return the created HTTP client", func() {
					result, err := subject.Authenticate(false)
					Expect(result).To(BeIdenticalTo(&httpClient))
					Expect(err).To(BeNil())
				})

				Context("When savecreds is true", func() {
					BeforeEach(func() {
						credsFile = new(storagefakes.FakeFile)
						credsCreateError = nil
					})

					It("should save the access token", func() {
						subject.Authenticate(true)
						Expect(fs.CreateCallCount()).To(Equal(1))
						Expect(fs.CreateArgsForCall(0)).To(Equal("flickr-credentials"))
						Expect(credsFile.WriteCallCount()).To(BeNumerically(">", 0))
					})
				})

				Context("When savecreds is false", func() {
					It("should not save the access token", func() {
						subject.Authenticate(false)
						Expect(fs.CreateCallCount()).To(Equal(0))
					})
				})
			})
		})
	})
})
