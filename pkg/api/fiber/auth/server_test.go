package auth_test

import (
	"bytes"
	"encoding/json"
	"github.com/arya-analytics/aryacore/pkg/api"
	"github.com/arya-analytics/aryacore/pkg/api/fiber/auth"
	baseauth "github.com/arya-analytics/aryacore/pkg/auth"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
)

func createUser(ds *mock.DataSourceMem) *models.User {
	pwd, err := baseauth.HashPassword("password")
	Expect(err).To(BeNil())
	user := &models.User{
		ID:       uuid.New(),
		Username: "GoodUser",
		Password: pwd,
	}
	Expect(ds.NewCreate().Model(user).Exec(ctx)).To(Succeed())
	return user
}

func marshallRequest(j interface{}) io.Reader {
	b, err := json.Marshal(j)
	Expect(err).To(BeNil())
	return bytes.NewReader(b)
}

var _ = Describe("Server", func() {
	var (
		app *fiber.App
		ds  *mock.DataSourceMem
	)
	BeforeEach(func() {
		app = fiber.New()
		ds = mock.NewDataSourceMem()
		svc := baseauth.NewService(ds.Exec)
		auth.NewServer(svc).BindTo(app)
	})
	Describe("Login", func() {
		It("Should return invalid credentials if the user isn't found", func() {
			createUser(ds)
			req := httptest.NewRequest("POST", "/auth/login", marshallRequest(map[string]string{
				"username": "BadUser",
				"password": "BadPassword",
			}))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
			var res map[string]interface{}
			unmarshalResponse(resp, &res)
			Expect(res["type"]).To(Equal(float64(api.ErrorTypeAuthentication)))
			Expect(res["message"]).To(Equal("Invalid credentials."))
		})
		It("Should return a good status and set a cookie if the login was successful", func() {
			user := createUser(ds)
			req := httptest.NewRequest("POST", "/auth/login", marshallRequest(map[string]string{
				"username": user.Username,
				"password": "password",
			}))
			req.Header.Set("Content-Type", "application/json")
			resp, err := app.Test(req)
			Expect(err).To(BeNil())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Cookies()).To(HaveLen(1))
			Expect(resp.Cookies()[0].Name).To(Equal("token"))
			var res map[string]interface{}
			unmarshalResponse(resp, &res)
			Expect(res["token"]).To(Equal(resp.Cookies()[0].Value))
		})
	})
})
