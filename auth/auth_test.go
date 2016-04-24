package auth_test

import (
	"errors"

	. "github.com/sgravrock/flickr-to-go-go/auth"

	"github.com/mrjones/oauth"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sgravrock/flickr-to-go-go/auth/authfakes"
)

var _ = Describe("Auth", func() {
	var oauthClient authfakes.FakeOauthClient
	var oauthConsumer *authfakes.FakeOauthConsumer
	var ui *authfakes.FakeUiAdapter

	BeforeEach(func() {
		oauthClient = *new(authfakes.FakeOauthClient)
		oauthConsumer = new(authfakes.FakeOauthConsumer)
		ui = new(authfakes.FakeUiAdapter)
		oauthClient.NewConsumerReturns(oauthConsumer)
	})

	Context("When obtaining a request token fails", func() {
		var tokenError error

		BeforeEach(func() {
			tokenError = errors.New("nope")
			oauthConsumer.GetRequestTokenAndUrlReturns(nil, "", tokenError)
		})

		It("should return the error", func() {
			t, err := Authenticate("theKey", "theSecret", &oauthClient, ui)
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
			Authenticate("theKey", "theSecret", &oauthClient, ui)
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
				t, err := Authenticate("theKey", "theSecret", &oauthClient, ui)
				Expect(t).To(BeNil())
				Expect(err).To(BeIdenticalTo(codeError))
			})
		})

		Context("When the user enters an access code", func() {
			BeforeEach(func() {
				ui.PromptForAccessCodeReturns("the code", nil)
			})

			It("should authorize the request token and code", func() {
				Authenticate("theKey", "theSecret", &oauthClient, ui)
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
					t, err := Authenticate("theKey", "theSecret", &oauthClient, ui)
					Expect(t).To(BeNil())
					Expect(err).To(BeIdenticalTo(authError))
				})

			})

			Context("When authorization succeeds", func() {
				var accessToken oauth.AccessToken

				BeforeEach(func() {
					accessToken = oauth.AccessToken{
						"access token", "", nil,
					}
					oauthConsumer.AuthorizeTokenReturns(&accessToken, nil)
				})

				It("should return the acesss token", func() {
					result, err := Authenticate("theKey", "theSecret", &oauthClient, ui)
					Expect(result).To(Equal(&accessToken))
					Expect(err).To(BeNil())
				})
			})
		})
	})
})
