package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Models", func() {
	Describe("Binding Values", func() {
		var c *storage.ChannelConfig
		var m *storage.Model
		BeforeEach(func() {
			c = &storage.ChannelConfig{}
			m = &storage.Model{Dest: c}
		})
		It("Should set the model values correctly", func() {
			id := uuid.New()
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
			id := uuid.New()
			c := &storage.ChannelConfig{Name: "Hello", ID: id}
			m := &storage.Model{Dest: c}
			mv := m.MapVals()
			Expect(mv).To(Equal(storage.ModelValues{"Name": "Hello", "ID": id}))
		})
	})
})
