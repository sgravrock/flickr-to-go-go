package auth

import (
	"encoding/gob"

	"github.com/mrjones/oauth"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

type UiAdapter interface {
	PromptForAccessCode(url string) (string, error)
}

func Authenticate(key string, secret string, filestore storage.Storage,
	savecreds bool, oauthClient OauthClient, ui UiAdapter) (*oauth.AccessToken, error) {

	f, err := filestore.Open("flickr-credentials")
	if err == nil {
		defer f.Close()
		result := new(oauth.AccessToken)
		err := gob.NewDecoder(f).Decode(result)
		return result, err
	}

	if oauthClient == nil {
		oauthClient = NewOauthClient()
	}
	if ui == nil {
		ui = NewConsoleUiAdapter()
	}

	accessToken, err := authLive(key, secret, oauthClient, ui)
	if err != nil || !savecreds {
		return accessToken, err
	}

	f, err = filestore.Create("flickr-credentials")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	err = gob.NewEncoder(f).Encode(accessToken)
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func authLive(key string, secret string, oauthClient OauthClient,
	ui UiAdapter) (*oauth.AccessToken, error) {

	consumer := oauthClient.NewConsumer(
		key,
		secret,
		"https://www.flickr.com/services/oauth/request_token",
		"https://www.flickr.com/services/oauth/authorize",
		"https://www.flickr.com/services/oauth/access_token")
	consumer.SetAdditionalParams(map[string]string{
		"perms": "read",
	})

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
