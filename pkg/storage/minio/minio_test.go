package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	mock2 "github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Minio Engine", func() {
	Describe("adapter", func() {
		var a storage.Adapter
		BeforeEach(func() {
			a = engine.NewAdapter()
		})
		Describe("New adapter", func() {
			It("Should create a new adapter without error", func() {
				Expect(reflect.TypeOf(a.ID())).To(Equal(reflect.TypeOf(uuid.New())))
			})
		})
		Describe("Is adapter", func() {
			Context("adapter is the correct type", func() {
				It("Should return true", func() {
					Expect(engine.IsAdapter(a)).To(BeTrue())
				})
			})
			Context("adapter is the incorrect type", func() {
				It("Should return false", func() {
					e := &mock.MDEngine{}
					ba := e.NewAdapter()
					Expect(engine.IsAdapter(ba)).To(BeFalse())
				})
			})
		})
		Context("Conn", func() {
			Describe("Binding an invalid adapter", func() {
				e := &mock.MDEngine{}
				ba := e.NewAdapter()
				Expect(func() {
					engine.NewRetrieve(ba)
				}).To(Panic())
			})
		})
	})
	Describe("Catalog", func() {
		Describe("Contains", func() {
			Context("Model in catalog", func() {
				It("Should return true", func() {
					Expect(engine.ShouldHandle(&models.ChannelChunkReplica{})).To(BeTrue())
				})
			})
			Context("Model not in catalog", func() {
				It("Should return false", func() {
					Expect(engine.ShouldHandle(&mock2.ModelB{})).To(BeFalse())
				})
			})
			Context("A model field that minio storage needs to handle not specified", func() {
				It("Should return false", func() {
					Expect(engine.ShouldHandle(&models.ChannelChunkReplica{}, "RangeReplicaID")).To(BeFalse())
				})
			})
			Context("A model field that minio needs to handle specified", func() {
				It("Should return true", func() {
					Expect(engine.ShouldHandle(&models.ChannelChunkReplica{}, "Telem")).To(BeTrue())
				})
			})
		})
	})
})
