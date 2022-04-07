package auth_test

import (
	"github.com/arya-analytics/aryacore/pkg/auth"
	"github.com/arya-analytics/aryacore/pkg/models"
	querymock "github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var (
		s  *auth.Service
		ds *querymock.DataSourceMem
	)
	BeforeEach(func() {
		ds = querymock.NewDataSourceMem()
		s = auth.NewService(ds.Exec)
	})
	Describe("Login", func() {
		It("Should log the user in correctly", func() {
			hash, err := auth.GenerateFromPassword("password")
			Expect(err).To(BeNil())
			user := &models.User{
				ID:       uuid.New(),
				Username: "root",
				Password: hash,
			}
			Expect(ds.NewCreate().Model(user).Exec(ctx)).To(Succeed())

			resUser, err := s.Login(ctx, user.Username, "password")
			Expect(err).To(BeNil())
			Expect(resUser.ID).To(Equal(user.ID))
		})
		It("Should return the correct auth error if the user isn't found", func() {
			_, err := s.Login(ctx, "root", "password")
			Expect(err.(auth.Error).Type).To(Equal(auth.ErrorTypeUserNotFound))
		})
		It("Should return the correct auth error if the credentials are invalid", func() {
			hash, err := auth.GenerateFromPassword("password")
			Expect(err).To(BeNil())
			user := &models.User{
				ID:       uuid.New(),
				Username: "root",
				Password: hash,
			}
			Expect(ds.NewCreate().Model(user).Exec(ctx)).To(Succeed())

			_, err = s.Login(ctx, user.Username, "wrongpassword")
			Expect(err.(auth.Error).Type).To(Equal(auth.ErrorTypeInvalidCredentials))
		})
	})
})
