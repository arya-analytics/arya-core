package model_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/arya-analytics/aryacore/pkg/util/model"
)

var _ = Describe("Catalog", func() {
	Context("Valid Catalog", func() {
		Context("Single Model", func() {
			It("Should return a new model of the correct type", func() {
				catalog := model.Catalog{
					&storage.ChannelConfig{},
				}
				Expect(model.NewReflect(catalog.New(&storage.ChannelConfig{})).
					Pointer()).To(Equal(&storage.ChannelConfig{}))
			})
		})
		Context("Chain of Models", func() {
			It("Should return a new chain of models of the correct type",
				func() {
					catalog := model.Catalog{
						&storage.ChannelConfig{},
					}
					Expect(model.NewReflect(catalog.New(&[]*storage.ChannelConfig{})).
						Pointer()).To(Equal(&[]*storage.ChannelConfig{}))
				})
		})
	})
	Describe("Contains", func() {
		Context("Catalog contains the model", func() {
			It("Should return true", func() {
				catalog := model.Catalog{
					&storage.ChannelConfig{},
				}
				Expect(catalog.Contains(&storage.ChannelConfig{})).To(BeTrue())
			})
		})
		Context("Catalog does not contain th model", func() {
			It("Should return false", func() {
				catalog := model.Catalog{
					&storage.ChannelConfig{},
				}
				Expect(catalog.Contains(&storage.ChannelChunk{})).To(BeFalse())
			})
		})
	})
	Context("Edge cases + errors", func() {
		It("Should panic when the catalog doesn't contain pointers", func() {
			catalog := model.Catalog{
				storage.ChannelConfig{},
			}
			Expect(func() {
				catalog.New(storage.ChannelConfig{})
			}).To(Panic())
		})
		It("Should panic when the model cannot be found in the catalog", func() {
			catalog := model.Catalog{
				&storage.ChannelConfig{},
			}
			Expect(func() {
				catalog.New(&storage.Node{})
			}).To(Panic())
		})
	})
})
