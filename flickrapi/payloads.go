package flickrapi

type FlickrPayload interface {
	Basics() *FlickrPayloadBasics
}

type FlickrPaginatedPayload interface {
	FlickrPayload
	PageInfo() *FlickrPayloadPageInfo
}

type FlickrPayloadBasics struct {
	Stat    string
	Message string
}

type FlickrPayloadPageInfo struct {
	Page  int
	Pages int
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

func (p *PeoplePhotosPayload) PageInfo() *FlickrPayloadPageInfo {
	return &(p.Photos.FlickrPayloadPageInfo)
}

type PeoplePhotosPhotos struct {
	FlickrPayloadPageInfo
	Photo []PhotoListEntry
}

type PhotoListEntry struct {
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
