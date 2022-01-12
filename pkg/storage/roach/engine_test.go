package roach_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	"github.com/arya-analytics/aryacore/pkg/storage/stub"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var dummyEngine = &roach.Engine{
	Driver: roach.DriverSQLite,
}

var _ = Describe("Engine", func() {
	Describe("Adapter", func() {
		Describe("New Adapter", func() {
			It("Should create a new adapter without error", func() {
				a := dummyEngine.NewAdapter()
				Expect(len(a.ID().String())).To(Equal(len(uuid.New().String())))
			})
		})
		Describe("Is Adapter", func() {
			Context("Adapter is the correct type", func() {
				It("Should return true", func() {
					a := dummyEngine.NewAdapter()
					Expect(dummyEngine.IsAdapter(a)).To(BeTrue())
				})
			})
			Context("Adapter is the incorrect type", func() {
				It("Should return false", func() {
					e := &stub.MDEngine{}
					a := e.NewAdapter()
					Expect(dummyEngine.IsAdapter(a)).To(BeFalse())
				})
			})
		})
	})
	Describe("Migrations", func() {
		Describe("Init Migrations", func() {
			It("Should initialize the migrations without error", func() {
				a := dummyEngine.NewAdapter()
				ctx := context.Background()
				err := dummyEngine.Migrate(ctx, a)
				Expect(err).To(BeNil())
			})
			It("Should create all of the tables correctly", func() {
				a := dummyEngine.NewAdapter()
				ctx := context.Background()
				_ = dummyEngine.Migrate(ctx, a)
				err := dummyEngine.VerifyMigrations(ctx, a)
				Expect(err).To(BeNil())
			})
		})
	})
})
