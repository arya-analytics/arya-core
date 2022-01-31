package minio_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryMigrate", func() {
	var err error
	BeforeEach(func() {
		err = mockEngine.NewMigrate(mockAdapter).Exec(mockCtx)
	})
	It("Should migrate without error", func() {
		Expect(err).To(BeNil())
	})
	It("Should execute all migrations successfully", func() {
		vErr := mockEngine.NewMigrate(mockAdapter).Verify(mockCtx)
		Expect(vErr).To(BeNil())
	})
})
