package model_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"

	"github.com/arya-analytics/aryacore/pkg/util/model"
)

var _ = Describe("Catalog", func() {
	Context("Valid Catalog", func() {
		Context("Single Model", func() {
			It("Should return a new model of the correct type", func() {
				catalog := model.Catalog{
					&models.ChannelConfig{},
				}
				Expect(model.NewReflect(catalog.New(&models.ChannelConfig{})).
					Pointer()).To(Equal(&models.ChannelConfig{}))
			})
		})
		Context("Chain of Models", func() {
			It("Should return a new chain of models of the correct type",
				func() {
					catalog := model.Catalog{
						&models.ChannelConfig{},
					}
					Expect(model.NewReflect(catalog.New(&[]*models.ChannelConfig{})).
						Pointer()).To(Equal(&[]*models.ChannelConfig{}))
				})
		})
	})
	Describe("Contains", func() {
		Context("Catalog contains the model", func() {
			It("Should return true", func() {
				catalog := model.Catalog{
					&models.ChannelConfig{},
				}
				Expect(catalog.Contains(&models.ChannelConfig{})).To(BeTrue())
			})
		})
		Context("Catalog does not contain th model", func() {
			It("Should return false", func() {
				catalog := model.Catalog{
					&models.ChannelConfig{},
				}
				Expect(catalog.Contains(&models.ChannelChunk{})).To(BeFalse())
			})
		})
	})
	Context("Edge cases + errors", func() {
		It("Should panic when the catalog doesn't contain pointers", func() {
			catalog := model.Catalog{
				models.ChannelConfig{},
			}
			Expect(func() {
				catalog.New(models.ChannelConfig{})
			}).To(Panic())
		})
		It("Should panic when the model cannot be found in the catalog", func() {
			catalog := model.Catalog{
				&models.ChannelConfig{},
			}
			Expect(func() {
				catalog.New(&models.Node{})
			}).To(Panic())
		})
	})
	Describe("Data Source", func() {
		It("Should retrieve the data source item correctly", func() {
			ds := model.NewDataSource()
			t := reflect.TypeOf(mock.ModelA{})
			ds.Write(model.NewReflect(&[]*mock.ModelA{}))
			m := ds.Retrieve(t)
			Expect(m.Type()).To(Equal(reflect.TypeOf(mock.ModelA{})))
		})
		It("Should create a new slice for a nonexistent item", func() {
			ds := model.NewDataSource()
			t := reflect.TypeOf(mock.ModelA{})
			m := ds.Retrieve(t)
			Expect(m.Type()).To(Equal(reflect.TypeOf(mock.ModelA{})))
		})
	})
})
