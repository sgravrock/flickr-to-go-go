package flickrapi

type FlickrPayload interface {
	Basics() *FlickrPayloadBasics
}

type FlickrPayloadBasics struct {
	Stat    string
	Message string
}

type TestLoginPayload struct {
	FlickrPayloadBasics
	User TestLoginUser
}

func (p *TestLoginPayload) Basics() *FlickrPayloadBasics {
	return &(p.FlickrPayloadBasics)
}

type TestLoginUser struct {
	Id       string
	Username TestLoginUsername
}

type TestLoginUsername struct {
	Content string `json:"_content"`
}
