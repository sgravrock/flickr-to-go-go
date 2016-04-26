package flickrapi_test

import (
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
		Context("with a successful response", func() {
			var payload []PhotoInfo
			var err error

			BeforeEach(func() {
				json := `{"photos":{"page":1,"pages":"1","perpage":2,
				"total":"2","photo":[{"id":"123","owner":"1234@N02",
				"secret":"asdf","server":"1518","farm":2,"title":"t1",
				"ispublic":1,"isfriend":0,"isfamily":0},{"id":"456",
				"owner":"1234@N02","secret":"qwer","server":"1521",
				"farm":2,"title":"t2","ispublic":0,"isfriend":1,"isfamily":0}
				]},"stat":"ok"}`
				ts.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET",
							"/services/rest/",
							"method=flickr.people.getPhotos&user_id=me&per_page=2&format=json&nojsoncallback=1"),
						ghttp.RespondWith(200, json),
					),
				)

				payload, err = subject.GetPhotos(2)
			})

			It("should populate the payload", func() {
				expected := []PhotoInfo{
					PhotoInfo{
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
					PhotoInfo{
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
				Expect(payload).To(Equal(expected))
			})

			It("should return nil", func() {
				Expect(err).To(BeNil())
			})
		})
	})
})
