package auth

import "github.com/mrjones/oauth"

type OauthClient interface {
	NewConsumer(key string, secret string, requestTokenUrl string,
		authorizeTokenUrl string, accessTokenUrl string) OauthConsumer
}

type OauthConsumer interface {
	GetRequestTokenAndUrl(callbackUrl string) (rtoken *oauth.RequestToken,
		loginUrl string, err error)
	AuthorizeToken(rtoken *oauth.RequestToken, verificationCode string) (atoken *oauth.AccessToken, err error)
	SetAdditionalParams(params map[string]string)
}

func NewOauthClient() OauthClient {
	return &defaultOauthClient{}
}

type defaultOauthClient struct{}

func (c *defaultOauthClient) NewConsumer(key string, secret string,
	requestTokenUrl string, authorizeTokenUrl string,
	accessTokenUrl string) OauthConsumer {

	wrappedConsumer := oauth.NewConsumer(
		key,
		secret,
		oauth.ServiceProvider{
			RequestTokenUrl:   requestTokenUrl,
			AuthorizeTokenUrl: authorizeTokenUrl,
			AccessTokenUrl:    accessTokenUrl,
		})
	return defaultOauthConsumer{wrappedConsumer}
}

type defaultOauthConsumer struct {
	wrappedConsumer *oauth.Consumer
}

func (c defaultOauthConsumer) GetRequestTokenAndUrl(callbackUrl string) (
	rtoken *oauth.RequestToken, loginUrl string, err error) {

	return c.wrappedConsumer.GetRequestTokenAndUrl(callbackUrl)
}

func (c defaultOauthConsumer) AuthorizeToken(rtoken *oauth.RequestToken,
	verificationCode string) (atoken *oauth.AccessToken, err error) {

	return c.wrappedConsumer.AuthorizeToken(rtoken, verificationCode)
}

func (c defaultOauthConsumer) SetAdditionalParams(params map[string]string) {
	c.wrappedConsumer.AdditionalAuthorizationUrlParams = params
}
