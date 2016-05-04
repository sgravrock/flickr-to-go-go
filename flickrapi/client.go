package flickrapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type Client interface {
	// Low-level interface
	Get(method string, params map[string]string, payload FlickrPayload) error

	// Higher-level interfaces for specific requests
	GetUsername() (string, error)
	GetPhotos(pageSize int) ([]PhotoInfo, error)
}

func NewClient(authenticatedHttpClient *http.Client, url string) Client {
	return flickrClient{authenticatedHttpClient, url}
}

type flickrClient struct {
	httpClient *http.Client
	url        string
}

func (c flickrClient) Get(method string, params map[string]string, payload FlickrPayload) error {
	url, err := c.buildUrl(method, params)
	if err != nil {
		return err
	}
	response, err := c.httpClient.Get(url)
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

func (c flickrClient) getPaged(method string, params map[string]string,
	payload FlickrPaginatedPayload, addPage func()) error {

	pagenum := 1

	for {
		params["page"] = strconv.Itoa(pagenum)
		err := c.Get(method, params, payload)
		if err != nil {
			return err
		}

		addPage()
		numPages := payload.PageInfo().Pages

		if numPages == 0 || pagenum >= numPages {
			return nil
		}

		pagenum++
	}
}

func (c flickrClient) buildUrl(method string, params map[string]string) (string, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return "", err
	}
	u.Path = "/services/rest/"

	q := u.Query()
	q.Set("method", method)
	q.Set("format", "json")
	q.Set("nojsoncallback", "1")
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
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

func (c flickrClient) GetPhotos(pageSize int) ([]PhotoInfo, error) {
	payload := PeoplePhotosPayload{}
	result := []PhotoInfo{}
	params := map[string]string{
		"user_id":  "me",
		"per_page": strconv.Itoa(pageSize),
	}
	err := c.getPaged("flickr.people.getPhotos", params, &payload, func() {
		result = append(result, payload.Photos.Photo...)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
