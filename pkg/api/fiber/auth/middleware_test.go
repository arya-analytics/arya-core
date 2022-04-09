package auth_test

import (
	"encoding/json"
	"github.com/arya-analytics/aryacore/pkg/api"
	"github.com/arya-analytics/aryacore/pkg/api/fiber/auth"
	baseauth "github.com/arya-analytics/aryacore/pkg/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func unmarshalResponse(resp *http.Response, into interface{}) {
	body, err := ioutil.ReadAll(resp.Body)
	Expect(err).To(BeNil())
	err = json.Unmarshal(body, into)
	if err != nil {
		log.Warn(string(body))
	}
	Expect(err).To(BeNil())
}

var _ = Describe("Middleware", func() {
	var app *fiber.App
	BeforeEach(func() {
		app = fiber.New()
	})
	Describe("TokenMiddleware", func() {
		var (
			req *http.Request
		)
		BeforeEach(func() {
			app.Use(auth.TokenMiddleware)
			app.Get("/hello", func(c *fiber.Ctx) error { return c.JSON(map[string]string{"Hello": "World!"}) })
			req = httptest.NewRequest("GET", "/hello", nil)
		})
		It("Should return an unauthorized error if no token is provided", func() {
			resp, err := app.Test(req)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
			var res map[string]interface{}
			unmarshalResponse(resp, &res)
			Expect(res["type"]).To(Equal(float64(api.ErrorTypeUnauthorized)))
			Expect(res["message"]).To(Equal("No authentication token provided. Please provide token as cookie or in headers."))
		})
		Context("Token as Cookie", func() {
			It("Should return an invalid token error if an invalid token is provided", func() {
				req.AddCookie(&http.Cookie{Name: "token", Value: "invalid"})
				resp, err := app.Test(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
				var res map[string]interface{}
				unmarshalResponse(resp, &res)
				Expect(res["type"]).To(Equal(float64(api.ErrorTypeUnauthorized)))
				Expect(res["message"]).To(Equal("Invalid authentication token provided."))
			})
			It("Should return a good response if a valid token is provided", func() {
				token, err := baseauth.NewToken(uuid.New())
				Expect(err).To(BeNil())
				req.AddCookie(&http.Cookie{Name: "token", Value: token})
				resp, err := app.Test(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
				var res map[string]interface{}
				unmarshalResponse(resp, &res)
				Expect(res["Hello"]).To(Equal("World!"))
			})
		})
		Context("Token as Header", func() {
			It("Should return an invalid token error if an invalid token is provided", func() {
				req.Header.Set("Authorization", "Bearer invalid")
				resp, err := app.Test(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
				var res map[string]interface{}
				unmarshalResponse(resp, &res)
				Expect(res["type"]).To(Equal(float64(api.ErrorTypeUnauthorized)))
				Expect(res["message"]).To(Equal("Invalid authentication token provided."))
			})
			It("Should return an invalid token error if the token is improperly formatter", func() {
				token, err := baseauth.NewToken(uuid.New())
				req.Header.Set("Authorization", token)
				Expect(err).To(BeNil())
				resp, err := app.Test(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
				var res map[string]interface{}
				unmarshalResponse(resp, &res)
				Expect(res["type"]).To(Equal(float64(api.ErrorTypeInvalidArguments)))
				Expect(res["message"]).To(Equal("Invalid authorization header. Expected format: 'Authorization: Bearer <token>'."))
			})
			It("Should return a good response if a valid token is provided", func() {
				token, err := baseauth.NewToken(uuid.New())
				req.Header.Set("Authorization", "Bearer "+token)
				Expect(err).To(BeNil())
				resp, err := app.Test(req)
				Expect(err).To(BeNil())
				var res map[string]interface{}
				unmarshalResponse(resp, &res)
				Expect(res["Hello"]).To(Equal("World!"))
				Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
			})

		})
	})
})
