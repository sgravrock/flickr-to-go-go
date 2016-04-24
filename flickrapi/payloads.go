package flickrapi

type TestLoginPayload struct {
	Stat string
	User TestLoginUser
}

type TestLoginUser struct {
	Id       string
	Username TestLoginUsername
}

type TestLoginUsername struct {
	Content string `json:"_content"`
}
