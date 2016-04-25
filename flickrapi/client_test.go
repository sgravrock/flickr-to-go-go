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
			var payload map[string]interface{}
			var err error

			BeforeEach(func() {
				payload = map[string]interface{}{}
				setupTestLoginSuccess(ts)
				payload, err = subject.Get("flickr.test.login", nil)
			})

			It("should deserialize JSON", func() {
				Expect(payload["stat"]).To(Equal("ok"))
				user, ok := payload["user"].(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(user["id"]).To(Equal("some-uid"))
				username, ok := user["username"].(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(username["_content"]).To(Equal("Some Username"))
			})

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

		Context("with a 200 response describing an error", func() {
			var payload map[string]interface{}
			var err error

			BeforeEach(func() {
				payload = map[string]interface{}{}
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
				payload, err = subject.Get("flickr.test.login", nil)
			})

			It("should return an error", func() {
				Expect(err).NotTo(BeNil())
				Expect(err.Error()).To(ContainSubstring("nope"))
			})
		})

		Context("with a non-200 response", func() {
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
				_, err = subject.Get("flickr.test.login", nil)
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

			It("should not return an error", func() {
				Expect(err).To(BeNil())
			})
		})

	})
})
