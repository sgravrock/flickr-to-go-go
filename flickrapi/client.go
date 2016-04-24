package flickrapi

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type Client interface {
	// Low-level interface
	Get(method string, params map[string]string, payload interface{}) error

	// Higher-level interfaces for specific requests
	GetUsername() (string, error)
}

func NewClient(authenticatedHttpClient *http.Client, url string) Client {
	return flickrClient{authenticatedHttpClient, url}
}

type flickrClient struct {
	httpClient *http.Client
	url        string
}

func (c flickrClient) Get(method string, params map[string]string, payload interface{}) error {
	// TODO: include params
	u, err := url.Parse(c.url)
	if err != nil {
		return err
	}
	u.Path = "/services/rest/"
	q := u.Query()
	q.Set("method", method)
	q.Set("format", "json")
	q.Set("nojsoncallback", "1")
	u.RawQuery = q.Encode()
	response, err := c.httpClient.Get(u.String())
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(payload)
}

func (c flickrClient) GetUsername() (string, error) {
	payload := TestLoginPayload{}
	err := c.Get("flickr.test.login", nil, &payload)
	if err != nil {
		return "", err
	}
	// TODO handle non-"ok" stat
	return payload.User.Username.Content, nil
}
