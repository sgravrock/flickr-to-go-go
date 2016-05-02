package auth

import (
	"encoding/gob"

	"net/http"

	"github.com/mrjones/oauth"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

type UiAdapter interface {
	PromptForAccessCode(url string) (string, error)
}

type Authenticator interface {
	Authenticate(savecreds bool) (*http.Client, error)
}

func NewAuthenticator(key string, secret string, filestore storage.Storage,
	oauthClient OauthClient, ui UiAdapter) Authenticator {

	if oauthClient == nil {
		oauthClient = NewOauthClient()
	}
	if ui == nil {
		ui = NewConsoleUiAdapter()
	}

	return &defaultAuthenticator{key, secret, filestore, oauthClient, ui}
}

type defaultAuthenticator struct {
	key         string
	secret      string
	filestore   storage.Storage
	oauthClient OauthClient
	ui          UiAdapter
}

func (a *defaultAuthenticator) Authenticate(savecreds bool) (*http.Client, error) {

	consumer := a.oauthClient.NewConsumer(
		a.key,
		a.secret,
		"https://www.flickr.com/services/oauth/request_token",
		"https://www.flickr.com/services/oauth/authorize",
		"https://www.flickr.com/services/oauth/access_token")
	consumer.SetAdditionalParams(map[string]string{
		"perms": "read",
	})

	accessToken, err := getAccessToken(a.key, a.secret, a.filestore, consumer, a.ui)
	if err != nil {
		return nil, err
	}

	if savecreds {
		err := saveCredentials(a.filestore, accessToken)
		if err != nil {
			return nil, err
		}
	}

	return consumer.MakeHttpClient(accessToken)
}

func getAccessToken(key string, secret string, filestore storage.Storage,
	consumer OauthConsumer, ui UiAdapter) (*oauth.AccessToken, error) {

	// Try to read saved credentials. Handle all errors as if the file
	// didn't exist.
	f, err := filestore.Open("flickr-credentials")
	if err == nil {
		defer f.Close()
		result := new(oauth.AccessToken)
		err := gob.NewDecoder(f).Decode(result)
		if err == nil {
			return result, nil
		}
	}

	return authLive(key, secret, consumer, ui)
}

func authLive(key string, secret string, consumer OauthConsumer,
	ui UiAdapter) (*oauth.AccessToken, error) {

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

func saveCredentials(filestore storage.Storage, token *oauth.AccessToken) error {
	f, err := filestore.Create("flickr-credentials")
	if err != nil {
		return err
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(token)
}
