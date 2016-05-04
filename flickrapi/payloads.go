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
	Username WrappedString
}

type WrappedString struct {
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

type PhotoInfoPayload struct {
	FlickrPayloadBasics
	Photo PhotoInfo
}

type PhotoInfo struct {
	Id                string           `json:"id"`
	Secret            string           `json:"secret"`
	Server            string           `json:"server"`
	Farm              int              `json:"farm"`
	Dateuploaded      string           `json:"dateuploaded"`
	Isfavorite        int              `json:"isfavorite"`
	License           int              `json:"license"`
	SafetyLevel       int              `json:"safety_level"`
	Rotation          int              `json:"rotation"`
	OriginalSecret    string           `json:"originalsecret"`
	OriginalFormat    string           `json:"originalformat"`
	Owner             PhotoOwner       `json:"owner"`
	Title             WrappedString    `json:"title"`
	Description       WrappedString    `json:"description"`
	Visibility        PhotoVisibility  `json:"visibility"`
	Dates             PhotoDates       `json:"dates"`
	Permissions       PhotoPermissions `json:"permissions"`
	Editability       PhotoEditability `json:"editability"`
	PublicEditability PhotoEditability `json:"publiceditability"`
	Usage             PhotoUsage       `json:"usage"`
	Comments          WrappedInt       `json:"comments"`
	Notes             PhotoNotes       `json:"notes"`
	People            PhotoHasPeople   `json:"people"`
	Tags              TagList          `json:"tags"`
	Urls              UrlList          `json:"urls"`
	Media             string           `json:"media"`
}

func (p *PhotoInfoPayload) Basics() *FlickrPayloadBasics {
	return &(p.FlickrPayloadBasics)
}

type PhotoOwner struct {
	Nsid       string `json:"nsid"`
	Username   string `json:"username"`
	RealName   string `json:"realname"`
	Location   string `json:"location"`
	IconServer int    `json:"iconserver"`
	IconFarm   int    `json:"iconfarm"`
	PathAlias  string `json:"path_alias"`
}

type PhotoVisibility struct {
	Ispublic int `json:"ispublic"`
	Isfriend int `json:"isfriend"`
	Isfamily int `json:"isfamily"`
}

type PhotoDates struct {
	Posted           string `json:"posted"`
	Taken            string `json:"taken"`
	TakenGranularity int    `json:"takengranularity"`
	TakenUnknown     int    `json:"takenunknown"`
	LastUpdate       string `json:"lastupdate"`
}

type PhotoPermissions struct {
	Comment     int `json:"permcomment"`
	AddMetadata int `json:"permaddmeta"`
}

type PhotoEditability struct {
	CanComment     int `json:"cancomment"`
	CanAddMetadata int `json:"canaddmeta"`
}

type PhotoUsage struct {
	CanDownload int `json:"candownload"`
	CanBlog     int `json:"canblog"`
	CanPrint    int `json:"canprint"`
	CanShare    int `json:"canshare"`
}

type WrappedInt struct {
	Content int `json:"_content"`
}

type PhotoNotes struct {
	Notes []PhotoNote `json:"note"`
}

type PhotoNote struct {
	Id         string `json:"id"`
	Author     string `json:"author"`
	AuthorName string `json:"authorname"`
	X          string `json:"x"`
	Y          string `json:"y"`
	W          int    `json:"w"`
	H          int    `json:"h"`
	Content    string `json:"_content"`
}

type PhotoHasPeople struct {
	HasPeople int `json:"haspeople"`
}

type TagList struct {
	Tags []Tag `json:"tag"`
}

type Tag struct {
	Id         string `json:"id"`
	Author     string `json:"author"`
	AuthorName string `json:"authorname"`
	Raw        string `json:"raw"`
	Content    string `json:"_content"`
	MachineTag int    `json:"machine_tag"`
}

type UrlList struct {
	Urls []Url `json:"url"`
}

type Url struct {
	Type    string `json:"type"`
	Content string `json:"_content"`
}
