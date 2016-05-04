package flickrapi_test

import (
	"fmt"

	. "github.com/sgravrock/flickr-to-go-go/flickrapi"

	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func setupTestLoginSuccess(ts *ghttp.Server) {
	json := "{\"user\": { \"id\": \"some-uid\", \"username\": { \"_content\": \"Some Username\" } }, \"stat\": \"ok\" }"
	ts.AppendHandlers(
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("GET",
				"/services/rest/",
				"method=flickr.test.login&format=json&nojsoncallback=1"),
			ghttp.RespondWith(200, json),
		),
	)
}

func setupPhotlistPages(ts *ghttp.Server, pages []string) {
	for i := 0; i < len(pages); i++ {
		params := fmt.Sprintf("method=flickr.people.getPhotos&user_id=me&page=%d&per_page=2&format=json&nojsoncallback=1", i+1)
		ts.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET",
					"/services/rest/",
					params),
				ghttp.RespondWith(200, pages[i]),
			),
		)
	}
}

var _ = Describe("flickrapi.Client", func() {
	var (
		subject Client
		ts      *ghttp.Server
	)

	BeforeEach(func() {
		ts = ghttp.NewServer()
		subject = NewClient(http.DefaultClient, ts.URL())
	})

	Describe("Get", func() {
		Context("with a successful response", func() {
			var payload TestLoginPayload
			var err error

			BeforeEach(func() {
				payload = TestLoginPayload{}
				setupTestLoginSuccess(ts)
				err = subject.Get("flickr.test.login", nil, &payload)
			})

			It("should deserialize the payload", func() {
				Expect(payload.Stat).To(Equal("ok"))
				Expect(payload.User.Id).To(Equal("some-uid"))
				Expect(payload.User.Username.Content).To(Equal("Some Username"))
			})

			It("should return nil", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("with a 200 response describing an error", func() {
			var payload TestLoginPayload
			var err error

			BeforeEach(func() {
				payload = TestLoginPayload{}
				err = nil
				json := "{\"stat\":\"fail\",\"code\":1,\"message\":\"nope\" }"
				ts.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET",
							"/services/rest/",
							"method=flickr.test.login&format=json&nojsoncallback=1"),
						ghttp.RespondWith(200, json),
					),
				)
				err = subject.Get("flickr.test.login", nil, &payload)
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("nope"))
			})
		})

		Context("with a non-200 response", func() {
			var payload TestLoginPayload
			var err error

			BeforeEach(func() {
				payload = TestLoginPayload{}
				ts.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET",
							"/services/rest/",
							"method=flickr.test.login&format=json&nojsoncallback=1"),
						ghttp.RespondWith(500, "oops"),
					),
				)
				err = subject.Get("flickr.test.login", nil, &payload)
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(Equal("flickr.test.login returned status 500"))
			})
		})
	})

	Describe("GetUsername", func() {
		Context("with a successful response", func() {
			var username string
			var err error

			BeforeEach(func() {
				setupTestLoginSuccess(ts)
				username, err = subject.GetUsername()
			})

			It("should return the username", func() {
				Expect(err).To(BeNil())
				Expect(username).To(Equal("Some Username"))
			})

			It("should return nil", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("GetPhotolist", func() {
		var result []PhotoListEntry
		var err error

		JustBeforeEach(func() {
			result, err = subject.GetPhotos(2)
		})

		Context("with a successful single-page response", func() {
			BeforeEach(func() {
				pages := []string{`{"photos":{"page":1,"pages":1,"perpage":2,
				"total":"2","photo":[{"id":"123","owner":"1234@N02",
				"secret":"asdf","server":"1518","farm":2,"title":"t1",
				"ispublic":1,"isfriend":0,"isfamily":0},{"id":"456",
				"owner":"1234@N02","secret":"qwer","server":"1521",
				"farm":2,"title":"t2","ispublic":0,"isfriend":1,"isfamily":0}
				]},"stat":"ok"}`}
				setupPhotlistPages(ts, pages)
			})

			It("should populate the payload", func() {
				expected := []PhotoListEntry{
					PhotoListEntry{
						Id:       "123",
						Owner:    "1234@N02",
						Secret:   "asdf",
						Server:   "1518",
						Farm:     2,
						Title:    "t1",
						Ispublic: 1,
						Isfriend: 0,
						Isfamily: 0,
					},
					PhotoListEntry{
						Id:       "456",
						Owner:    "1234@N02",
						Secret:   "qwer",
						Server:   "1521",
						Farm:     2,
						Title:    "t2",
						Ispublic: 0,
						Isfriend: 1,
						Isfamily: 0,
					},
				}
				Expect(result).To(Equal(expected))
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("with multiple pages", func() {
			Context("with an error", func() {
				BeforeEach(func() {
					ts.AppendHandlers(
						ghttp.CombineHandlers(
							ghttp.VerifyRequest("GET",
								"/services/rest/",
								"method=flickr.people.getPhotos&user_id=me&page=1&per_page=2&format=json&nojsoncallback=1"),
							ghttp.RespondWith(500, "oops"),
						),
					)
				})

				It("should return a nil photo list", func() {
					Expect(result).To(BeNil())
				})

				It("should return an error", func() {
					Expect(err).NotTo(BeNil())
				})
			})

			Context("with all successful responses", func() {
				BeforeEach(func() {
					pages := []string{
						`{"photos":{"page":1,"pages":2,"perpage":2,
						"total":"3","photo":[{"id":"123","owner":"1234@N02",
						"secret":"asdf","server":"1518","farm":2,"title":"t1",
						"ispublic":1,"isfriend":0,"isfamily":0},{"id":"456",
						"owner":"1234@N02","secret":"qwer","server":"1521",
						"farm":2,"title":"t2","ispublic":0,"isfriend":1,"isfamily":0}
						]},"stat":"ok"}`,
						`{"photos":{"page":2,"pages":2,"perpage":2,
						"total":"3","photo":[{"id":"789","owner":"1234@N02",
						"secret":"zxcv","server":"1518","farm":2,"title":"t3",
						"ispublic":0,"isfriend":0,"isfamily":1}
						]},"stat":"ok"}`,
					}
					setupPhotlistPages(ts, pages)
				})

				It("should return the list of photos", func() {
					expected := []PhotoListEntry{
						PhotoListEntry{
							Id:       "123",
							Owner:    "1234@N02",
							Secret:   "asdf",
							Server:   "1518",
							Farm:     2,
							Title:    "t1",
							Ispublic: 1,
							Isfriend: 0,
							Isfamily: 0,
						},
						PhotoListEntry{
							Id:       "456",
							Owner:    "1234@N02",
							Secret:   "qwer",
							Server:   "1521",
							Farm:     2,
							Title:    "t2",
							Ispublic: 0,
							Isfriend: 1,
							Isfamily: 0,
						},
						PhotoListEntry{
							Id:       "789",
							Owner:    "1234@N02",
							Secret:   "zxcv",
							Server:   "1518",
							Farm:     2,
							Title:    "t3",
							Ispublic: 0,
							Isfriend: 0,
							Isfamily: 1,
						},
					}
					Expect(result).To(Equal(expected))
				})

				It("should not return an error", func() {
					Expect(err).To(BeNil())
				})
			})
		})
	})

	Describe("GetPhotoInfo", func() {
		Context("When the request succeeds", func() {
			var result PhotoInfo
			var err error

			BeforeEach(func() {
				json := `{ "photo": { "id": "12345", "secret": "s3kr1t", "server": "2546", "farm": 3, "dateuploaded": "1256964640", "isfavorite": 1, "license": 5, "safety_level": 1, "rotation": 1, "originalsecret": "origsecret", "originalformat": "jpg",
				    "owner": { "nsid": "uid@N02", "username": "username", "realname": "realname", "location": "there", "iconserver": 2, "iconfarm": 3, "path_alias": "pa" },
				    "title": { "_content": "Title" },
				    "description": { "_content": "description" },
				    "visibility": { "ispublic": 1, "isfriend": 1, "isfamily": 1 },
				    "dates": { "posted": "1256964640", "taken": "2009-10-29 17:20:25", "takengranularity": 1, "takenunknown": 2, "lastupdate": "1257010788" },
				    "permissions": { "permcomment": 3, "permaddmeta": 2 }, "views": 33,
				    "editability": { "cancomment": 1, "canaddmeta": 1 },
				    "publiceditability": { "cancomment": 1, "canaddmeta": 1 },
				    "usage": { "candownload": 1, "canblog": 1, "canprint": 1, "canshare": 1 },
				    "comments": { "_content": 1 },
				    "notes": {
				      "note": [
				        { "id": "nid", "author": "uid@N02", "authorname": "username", "x": "135", "y": "310", "w": 69, "h": 52, "_content": "A note" }
				      ] },
				    "people": { "haspeople": 1 },
				    "tags": {
				      "tag": [
				        { "id": "tag ID", "author": "uid@N02", "authorname": "username", "raw": "The Tag", "_content": "thetag", "machine_tag": 1 }
				      ] },
				    "urls": {
				      "url": [
				        { "type": "photopage", "_content": "https://www.flickr.com/photos/uid@N02/12345/" }
				      ] }, "media": "photo" }, "stat": "ok" }`
				ts.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET",
							"/services/rest/",
							"method=flickr.photos.getInfo&format=json&nojsoncallback=1&photo_id=12345"),
						ghttp.RespondWith(200, json),
					),
				)

				result, err = subject.GetPhotoInfo("12345")
			})

			It("should return the photo info", func() {
				expected := PhotoInfo{
					Id:             "12345",
					Secret:         "s3kr1t",
					Server:         "2546",
					Farm:           3,
					Dateuploaded:   "1256964640",
					Isfavorite:     1,
					License:        5,
					SafetyLevel:    1,
					Rotation:       1,
					OriginalSecret: "origsecret",
					OriginalFormat: "jpg",
					Owner: PhotoOwner{
						Nsid:       "uid@N02",
						Username:   "username",
						RealName:   "realname",
						Location:   "there",
						IconServer: 2,
						IconFarm:   3,
						PathAlias:  "pa",
					},
					Title:       WrappedString{"Title"},
					Description: WrappedString{"description"},
					Visibility: PhotoVisibility{
						Ispublic: 1,
						Isfriend: 1,
						Isfamily: 1,
					},
					Dates: PhotoDates{
						Posted:           "1256964640",
						Taken:            "2009-10-29 17:20:25",
						TakenGranularity: 1,
						TakenUnknown:     2,
						LastUpdate:       "1257010788",
					},
					Permissions: PhotoPermissions{
						Comment:     3,
						AddMetadata: 2,
					},
					Editability: PhotoEditability{
						CanComment:     1,
						CanAddMetadata: 1,
					},
					PublicEditability: PhotoEditability{
						CanComment:     1,
						CanAddMetadata: 1,
					},
					Usage: PhotoUsage{
						CanDownload: 1,
						CanBlog:     1,
						CanPrint:    1,
						CanShare:    1,
					},
					Comments: WrappedInt{1},
					Notes: PhotoNotes{
						Notes: []PhotoNote{
							PhotoNote{
								Id:         "nid",
								Author:     "uid@N02",
								AuthorName: "username",
								X:          "135",
								Y:          "310",
								W:          69,
								H:          52,
								Content:    "A note",
							},
						},
					},
					People: PhotoHasPeople{1},
					Tags: TagList{
						Tags: []Tag{
							Tag{
								Id:         "tag ID",
								Author:     "uid@N02",
								AuthorName: "username",
								Raw:        "The Tag",
								Content:    "thetag",
								MachineTag: 1,
							},
						},
					},
					Urls: UrlList{
						Urls: []Url{
							Url{
								Type:    "photopage",
								Content: "https://www.flickr.com/photos/uid@N02/12345/",
							},
						},
					},
					Media: "photo",
				}
				Expect(err).To(BeNil())
				Expect(result).To(Equal(expected))
			})
		})
	})
})
