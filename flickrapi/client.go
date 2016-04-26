package flickrapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Client interface {
	// Low-level interface
	Get(method string, params map[string]string, payload FlickrPayload) error

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

func (c flickrClient) Get(method string, params map[string]string, payload FlickrPayload) error {
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
	if response.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("%s returned status %d", method, response.StatusCode)
		return errors.New(msg)
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(payload)
	if err != nil {
		return err
	}
	return verifyResponse(method, payload)
}

func verifyResponse(method string, payload FlickrPayload) error {
	basics := payload.Basics()
	if basics.Stat != "ok" {
		msg := fmt.Sprintf("%s failed: status %s, message %s",
			method, basics.Stat, basics.Message)
		return errors.New(msg)
	}
	return nil
}

func (c flickrClient) GetUsername() (string, error) {
	payload := TestLoginPayload{}
	err := c.Get("flickr.test.login", nil, &payload)
	if err != nil {
		return "", err
	}
	return payload.User.Username.Content, nil
}
