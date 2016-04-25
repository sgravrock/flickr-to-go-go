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
	Get(method string, params map[string]string) (map[string]interface{}, error)

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

func (c flickrClient) Get(method string, params map[string]string) (map[string]interface{}, error) {
	var payload map[string]interface{}
	// TODO: include params
	u, err := url.Parse(c.url)
	if err != nil {
		return payload, err
	}
	u.Path = "/services/rest/"
	q := u.Query()
	q.Set("method", method)
	q.Set("format", "json")
	q.Set("nojsoncallback", "1")
	u.RawQuery = q.Encode()
	response, err := c.httpClient.Get(u.String())
	if err != nil {
		return payload, err
	}
	if response.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("%s returned status %d", method, response.StatusCode)
		return payload, errors.New(msg)
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return payload, err
	}

	return payload, checkResponse(payload)
}

func checkResponse(payload map[string]interface{}) error {
	if payload["stat"] == "ok" {
		return nil
	}

	msg := fmt.Sprintf("API call failed with status: %s, message: %s",
		payload["stat"], payload["message"])
	return errors.New(msg)
}

func (c flickrClient) GetUsername() (string, error) {
	payload, err := c.Get("flickr.test.login", nil)
	if err != nil {
		return "", err
	}

	user, ok := payload["user"].(map[string]interface{})
	if !ok {
		return "", errors.New("Unexpected API call result format")
	}
	username, ok := user["username"].(map[string]interface{})
	if !ok {
		return "", errors.New("Unexpected API call result format")
	}
	content, ok := username["_content"].(string)
	if !ok {
		return "", errors.New("Unexpected API call result format")
	}
	return content, nil
}
