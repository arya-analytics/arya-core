package minio_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryMigrate", func() {
	var err error
	BeforeEach(func() {
		err = engine.NewMigrate().Exec(ctx)
	})
	It("Should migrate without error", func() {
		Expect(err).To(BeNil())
	})
	It("Should execute all migrations successfully", func() {
		vErr := engine.NewMigrate().Verify().Exec(ctx)
		Expect(vErr).To(BeNil())
	})
})
