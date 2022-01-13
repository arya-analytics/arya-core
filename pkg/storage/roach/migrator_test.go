package roach_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Migrator", func() {
	var mErr error
	var a storage.Adapter
	BeforeEach(func() {
		a = dummyEngine.NewAdapter()
		ctx := context.Background()
		mErr = dummyEngine.Migrate(ctx, a)
	})
	Describe("Init Migrations", func() {
		It("Should initialize the migrations without error", func() {
			Expect(mErr).To(BeNil())
		})
		It("Should create all of the tables correctly", func() {
			ctx := context.Background()
			err := dummyEngine.VerifyMigrations(ctx, a)
			Expect(err).To(BeNil())
		})
	})
})
