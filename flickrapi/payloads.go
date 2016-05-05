package flickrapi

import "errors"

type PhotoListEntry struct {
	Data map[string]interface{}
}

func (e *PhotoListEntry) Id() (string, error) {
	id, ok := e.Data["id"].(string)
	if !ok {
		return "", errors.New("Unexpected API result format (no id in photo list entry)")
	}

	return id, nil
}

func (e *PhotoListEntry) OriginalUrl() (string, error) {
	return requireString(e.Data, []string{"url_o"})
}
