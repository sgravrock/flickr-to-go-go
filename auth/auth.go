package auth

import (
	"encoding/gob"

	"fmt"

	"github.com/mrjones/oauth"
	"github.com/sgravrock/flickr-to-go-go/storage"
)

type UiAdapter interface {
	PromptForAccessCode(url string) (string, error)
}

func Authenticate(key string, secret string, filestore storage.Storage,
	savecreds bool, oauthClient OauthClient, ui UiAdapter) (*oauth.AccessToken, error) {

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

	accessToken, err := consumer.AuthorizeToken(requestToken, accessCode)
	if err != nil {
		return nil, err
	}

	if !savecreds {
		fmt.Println("Not saving creds")
		return accessToken, nil
	}

	fmt.Println("Saving creds")
	f, err := filestore.Create("flickr-credentials")
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
