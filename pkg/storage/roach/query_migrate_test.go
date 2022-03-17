package roach_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migrator", func() {
	Describe("Verify Migrations", func() {
		It("Should create all of the tables correctly", func() {
			err := engine.NewMigrate().Verify().Exec(ctx)
			Expect(err).To(BeNil())
		})
	})
})
