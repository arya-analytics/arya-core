package roach_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migrator", func() {
	BeforeEach(migrate)
	Describe("Init Migrations", func() {
		It("Should create all of the tables correctly", func() {
			err := mockEngine.NewMigrate(mockAdapter).Verify(mockCtx)
			Expect(err).To(BeNil())
		})
	})
})
