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
}

func NewOauthClient() OauthClient {
	return &defaultOauthClient{}
}

type defaultOauthClient struct {
}

func (c *defaultOauthClient) NewConsumer(key string, secret string,
	requestTokenUrl string, authorizeTokenUrl string,
	accessTokenUrl string) OauthConsumer {

	return oauth.NewConsumer(
		key,
		secret,
		oauth.ServiceProvider{
			RequestTokenUrl:   requestTokenUrl,
			AuthorizeTokenUrl: authorizeTokenUrl,
			AccessTokenUrl:    accessTokenUrl,
		})
}
