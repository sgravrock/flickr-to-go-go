package auth

import "github.com/mrjones/oauth"

type UiAdapter interface {
	PromptForAccessCode(url string) (string, error)
}

func Authenticate(key string, secret string, oauthClient OauthClient,
	ui UiAdapter) (*oauth.AccessToken, error) {
	if oauthClient == nil {
		oauthClient = NewOauthClient()
	}
	if ui == nil {
		ui = NewConsoleUiAdapter()
	}

	consumer := oauthClient.NewConsumer(
		key,
		secret,
		"https://www.flickr.com/services/oauth/request_token",
		"https://www.flickr.com/services/oauth/authorize",
		"https://www.flickr.com/services/oauth/access_token")

	requestToken, url, err := consumer.GetRequestTokenAndUrl("oob")
	if err != nil {
		return nil, err
	}

	accessCode, err := ui.PromptForAccessCode(url)
	if err != nil {
		return nil, err
	}

	return consumer.AuthorizeToken(requestToken, accessCode)
}
