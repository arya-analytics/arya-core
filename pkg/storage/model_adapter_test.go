package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Model Adapter", func() {
	Describe("SingleModelAdapter", func() {
		Describe("Binding Values", func() {
			var c *storage.ChannelConfig
			var m *storage.AdaptedModel
			BeforeEach(func() {
				c = &storage.ChannelConfig{}
				m = storage.NewAdaptedModel(c)
			})
			It("Should set the model values correctly", func() {
				var id = 445
				err := m.BindVals(storage.ModelValues{"Name": "Hello", "ID": id})
				Expect(err).To(BeNil())
				Expect(c.Name).To(Equal("Hello"))
				Expect(c.ID).To(Equal(id))
			})
			It("Should return an error if a non-existent field is provided", func() {
				err := m.BindVals(storage.ModelValues{"InvalidKey": "Invalid Value"})
				Expect(err).ToNot(BeNil())
			})
			It("Should return an error if an invalid type is provided", func() {
				err := m.BindVals(storage.ModelValues{"Name": 221})
				Expect(err).ToNot(BeNil())
			})
		})
		Describe("Mapping Values", func() {
			It("Should map all values correctly", func() {
				var id = 445
				c := &storage.ChannelConfig{Name: "Hello", ID: id}
				m := storage.NewAdaptedModel(c)
				mv := m.MapVals()
				Expect(mv).To(Equal(storage.ModelValues{"Name": "Hello", "ID": id}))
			})
		})
		Describe("Re-binding to interface", func() {
		})

	})
	FDescribe("Multi Model Adapter", func() {
		It("Should bind the source and des models correctly", func() {
			var destModels = []*storage.ChannelConfig{
				dummyModel,
			}
			var sourceModels []*storage.ChannelConfig
			ma := storage.NewModelAdapter(&sourceModels, &destModels)
			if err := ma.ExchangeToSource(); err != nil {
				log.Fatalln(err)
			}
		})
	})
})
