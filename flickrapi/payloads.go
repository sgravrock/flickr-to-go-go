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

type PeoplePhotosPayload struct {
	FlickrPayloadBasics
	Photos PeoplePhotosPhotos
}

func (p *PeoplePhotosPayload) Basics() *FlickrPayloadBasics {
	return &(p.FlickrPayloadBasics)
}

type PeoplePhotosPhotos struct {
	Photo []PhotoInfo
}

type PhotoInfo struct {
	Id       string
	Owner    string
	Secret   string
	Server   string
	Farm     int
	Title    string
	Ispublic int
	Isfriend int
	Isfamily int
}
