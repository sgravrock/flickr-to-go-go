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

func setupRecentPhotosPages(ts *ghttp.Server, pages []string, timestamp uint32) {
	for i := 0; i < len(pages); i++ {
		params := fmt.Sprintf("method=flickr.photos.recentlyUpdated&min_date=%d&page=%d&per_page=2&extras=url_o&format=json&nojsoncallback=1",
			timestamp, i+1)
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

func setupPhotolistPages(ts *ghttp.Server, pages []string) {
	for i := 0; i < len(pages); i++ {
		params := fmt.Sprintf("method=flickr.people.getPhotos&user_id=me&page=%d&per_page=2&extras=url_o&format=json&nojsoncallback=1", i+1)
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
			var result map[string]interface{}
			var err error

			BeforeEach(func() {
				setupTestLoginSuccess(ts)
				result, err = subject.Get("flickr.test.login", nil)
			})

			It("should deserialize the payload", func() {
				Expect(result["stat"]).To(Equal("ok"))
				user, ok := result["user"].(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(user["id"]).To(Equal("some-uid"))
				username, ok := user["username"].(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(username["_content"]).To(Equal("Some Username"))
			})

			It("should return nil", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("with a 200 response describing an error", func() {
			var result map[string]interface{}
			var err error

			BeforeEach(func() {
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
				result, err = subject.Get("flickr.test.login", nil)
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("nope"))
			})
		})

		Context("with a non-200 response", func() {
			var result map[string]interface{}
			var err error

			BeforeEach(func() {
				ts.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET",
							"/services/rest/",
							"method=flickr.test.login&format=json&nojsoncallback=1"),
						ghttp.RespondWith(500, "oops"),
					),
				)
				result, err = subject.Get("flickr.test.login", nil)
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

	Describe("GetRecentPhotoIds", func() {
		var result []string
		var err error

		JustBeforeEach(func() {
			result, err = subject.GetRecentPhotoIds(12345, 2)
		})

		Context("with a successful single-page response", func() {
			BeforeEach(func() {
				pages := []string{`{ "photos": { "page": 1, "pages": 1, "perpage": 2, "total": 2, 
    				"photo": [
      					{ "id": "11111"},
      					{ "id": "22222"}
					] }, "stat": "ok" }`}
				setupRecentPhotosPages(ts, pages, 12345)
			})

			It("should return the photo IDs", func() {
				Expect(result).To(Equal([]string{"11111", "22222"}))
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("with an error", func() {
			BeforeEach(func() {
				ts.AppendHandlers(
					ghttp.CombineHandlers(

						ghttp.VerifyRequest("GET",
							"/services/rest/",
							"method=flickr.photos.recentlyUpdated&min_date=12345&page=1&per_page=2&extras=url_o&format=json&nojsoncallback=1"),
						ghttp.RespondWith(500, "oops"),
					),
				)
			})

			It("should return a nil list", func() {
				Expect(result).To(BeNil())
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
			})
		})

		Context("with all successful responses", func() {
			BeforeEach(func() {
				pages := []string{
					`{ "photos": { "page": 1, "pages": 2, "perpage": 100, "total": 3, 
    				"photo": [
      					{ "id": "11111"},
      					{ "id": "22222"}
					] }, "stat": "ok" }`,
					`{ "photos": { "page": 2, "pages": 2, "perpage": 100, "total": 3, 
    				"photo": [
      					{ "id": "33333"}
					] }, "stat": "ok" }`,
				}
				setupRecentPhotosPages(ts, pages, 12345)
			})

			It("should return the photo IDs", func() {
				Expect(result).To(Equal([]string{"11111", "22222", "33333"}))
			})

			It("should not return an error", func() {
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
				setupPhotolistPages(ts, pages)
			})

			It("should populate the payload", func() {
				Expect(len(result)).To(Equal(2))
				Expect(result[0].Id()).To(Equal("123"))
				Expect(result[1].Id()).To(Equal("456"))
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("with an error", func() {
			BeforeEach(func() {
				ts.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET",
							"/services/rest/",
							"method=flickr.people.getPhotos&user_id=me&page=1&per_page=2&extras=url_o&format=json&nojsoncallback=1"),
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
				setupPhotolistPages(ts, pages)
			})

			It("should return the list of photos", func() {
				Expect(len(result)).To(Equal(3))
				Expect(result[0].Id()).To(Equal("123"))
				Expect(result[1].Id()).To(Equal("456"))
				Expect(result[2].Id()).To(Equal("789"))
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("GetPhotoInfo", func() {
		Context("When the request succeeds", func() {
			var result map[string]interface{}
			var err error

			BeforeEach(func() {
				json := `{ "stat": "ok", "photo": { "id": "12345", "title": { "_content": "Title" }}}`
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
				expected := map[string]interface{}{
					"id": "12345",
					"title": map[string]interface{}{
						"_content": "Title",
					},
				}

				Expect(err).To(BeNil())
				Expect(result).To(Equal(expected))
			})
		})
	})
})
