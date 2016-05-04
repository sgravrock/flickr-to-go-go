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
	Get(method string, params map[string]string) (map[string]interface{}, error)

	// Higher-level interfaces for specific requests
	GetUsername() (string, error)
	GetPhotos(pageSize int) ([]PhotoListEntry, error)
	GetPhotoInfo(photoId string) (map[string]interface{}, error)
}

func NewClient(authenticatedHttpClient *http.Client, url string) Client {
	return flickrClient{authenticatedHttpClient, url}
}

type flickrClient struct {
	httpClient *http.Client
	url        string
}

func (c flickrClient) Get(method string, params map[string]string) (map[string]interface{}, error) {
	url, err := c.buildUrl(method, params)
	if err != nil {
		return nil, err
	}
	response, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("%s returned status %d", method, response.StatusCode)
		return nil, errors.New(msg)
	}
	defer response.Body.Close()
	var payload map[string]interface{}
	err = json.NewDecoder(response.Body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	err = verifyResponse(method, payload)
	return payload, err
}

func (c flickrClient) getPaged(method string, params map[string]string,
	pageInfoKey string, addPage func(map[string]interface{}) error) error {

	pagenum := 1

	for {
		params["page"] = strconv.Itoa(pagenum)
		payload, err := c.Get(method, params)
		if err != nil {
			return err
		}

		err = addPage(payload)
		if err != nil {
			return err
		}

		pageInfo, ok := payload[pageInfoKey].(map[string]interface{})
		if !ok {
			msg := fmt.Sprintf("Unexpected API call result format (no %s)", pageInfoKey)
			return errors.New(msg)
		}
		n, ok := pageInfo["pages"].(float64)
		if !ok {
			msg := fmt.Sprintf("Unexpected API call result format (no %s.pages)", pageInfoKey)
			return errors.New(msg)
		}
		numPages := int(n)

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

func verifyResponse(method string, payload map[string]interface{}) error {
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
		return "", errors.New("Unexpected API call result format (no user)")
	}
	username, ok := user["username"].(map[string]interface{})
	if !ok {
		return "", errors.New("Unexpected API call result format (no user.username)")
	}
	value, ok := username["_content"].(string)
	if !ok {
		return "", errors.New("Unexpected API call result format (no user.username._content)")
	}
	return value, nil
}

func (c flickrClient) GetPhotos(pageSize int) ([]PhotoListEntry, error) {
	result := []PhotoListEntry{}
	params := map[string]string{
		"user_id":  "me",
		"per_page": strconv.Itoa(pageSize),
	}
	err := c.getPaged("flickr.people.getPhotos", params, "photos",
		func(pagePayload map[string]interface{}) error {
			wpr, ok := pagePayload["photos"].(map[string]interface{})
			if !ok {
				return errors.New("Unexpected API call result format (no photos)")
			}
			photos, ok := wpr["photo"].([]interface{})
			if !ok {
				return errors.New("Unexpected API call result format (no photos.photo)")
			}
			for _, p := range photos {
				photo, ok := p.(map[string]interface{})
				if !ok {
					return errors.New("Unexpected API call result format (non-object in photos.photo)")
				}
				result = append(result, PhotoListEntry{photo})
			}
			return nil
		})

	if err != nil {
		return nil, err
	}
	return result, nil

}

func (c flickrClient) GetPhotoInfo(photoId string) (map[string]interface{}, error) {
	params := map[string]string{"photo_id": photoId}
	// TODO: test and handle errors
	payload, _ := c.Get("flickr.photos.getInfo", params)
	photo, ok := payload["photo"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Unexpected API call result format (no photo)")
	}
	return photo, nil
}
