package auth_test

import (
	"github.com/arya-analytics/aryacore/pkg/auth"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Token", func() {
	Describe("Generate", func() {
		It("should generate a token", func() {
			token, err := auth.NewToken(uuid.New())
			Expect(err).ToNot(HaveOccurred())
			Expect(token).ToNot(BeEmpty())
		})
	})
	Describe("Validate", func() {
		It("Should return a nil error for a valid token", func() {
			token, err := auth.NewToken(uuid.New())
			Expect(err).ToNot(HaveOccurred())
			err = auth.ValidateToken(token)
			Expect(err).ToNot(HaveOccurred())
		})
		It("Should return an invalid credentials error for an invalid token", func() {
			err := auth.ValidateToken("invalid")
			Expect(err).To(HaveOccurred())
			authErr := err.(auth.Error)
			Expect(authErr.Type).To(Equal(auth.ErrorTypeInvalidCredentials))
			Expect(authErr.Error()).To(Equal("ErrorTypeInvalidCredentials - Invalid token"))
		})
	})

})
