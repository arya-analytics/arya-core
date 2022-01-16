package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Model Adapter", func() {
	Context("Single Model Adaptation", func() {
		Context("Between two models of the same type", func() {
			Context("No nested model", func() {
				var source *storage.ChannelConfig
				var dest *storage.ChannelConfig
				BeforeEach(func() {
					source = &storage.ChannelConfig{
						ID:   435,
						Name: "Cool Name",
					}
					dest = &storage.ChannelConfig{}
				})
				It("Should exchange to source", func() {
					opts := &storage.ModelAdapterOpts{
						Source:        dest,
						Dest:          source,
						CatalogSource: storage.Catalog(),
						CatalogDest:   storage.Catalog(),
					}
					ma := storage.NewModelAdapter(opts)
					err := ma.ExchangeToSource()
					Expect(err).To(BeNil())
					Expect(source.ID).To(Equal(435))
					Expect(source.ID).To(Equal(dest.ID))
					Expect(source.Name).To(Equal(dest.Name))
				})
				It("Should exchange to dest", func() {
					opts := &storage.ModelAdapterOpts{
						Source:        source,
						Dest:          dest,
						CatalogDest:   storage.Catalog(),
						CatalogSource: storage.Catalog(),
					}
					ma := storage.NewModelAdapter(opts)
					err := ma.ExchangeToDest()
					Expect(err).To(BeNil())
					Expect(source.ID).To(Equal(435))
					Expect(source.ID).To(Equal(dest.ID))
					Expect(source.Name).To(Equal(dest.Name))
				})
				It("Shouldn't create references between source and dest models", func() {
					opts := &storage.ModelAdapterOpts{
						Source:        source,
						Dest:          dest,
						CatalogSource: storage.Catalog(),
						CatalogDest:   storage.Catalog(),
					}
					ma := storage.NewModelAdapter(opts)
					if err := ma.ExchangeToDest(); err != nil {
						log.Fatalln(err)
					}
					source.Name = "Hello"
					Expect(dest.Name).To(Equal("Cool Name"))
				})
			})
			Context("With nested model", func() {
				var source *storage.ChannelConfig
				var dest *storage.ChannelConfig
				var node *storage.Node
				BeforeEach(func() {
					node = &storage.Node{
						ID: 24,
					}
					source = &storage.ChannelConfig{
						ID:   420,
						Name: "Even Cooler Name",
						Node: node,
					}
					dest = &storage.ChannelConfig{}
				})
				It("Should exchange to source", func() {
					opts := &storage.ModelAdapterOpts{
						Source: dest,
						Dest:   source,
					}
					ma := storage.NewModelAdapter(opts)
					err := ma.ExchangeToSource()
					Expect(err).To(BeNil())
					Expect(source.Node.ID).To(Equal(24))
					Expect(dest.Node.ID).To(Equal(source.Node.ID))
				})
				It("Should exchange to dest", func() {
					opts := &storage.ModelAdapterOpts{
						Source: source,
						Dest:   dest,
					}
					ma := storage.NewModelAdapter(opts)
					err := ma.ExchangeToDest()
					Expect(err).To(BeNil())
					Expect(source.Node.ID).To(Equal(24))
					Expect(dest.Node.ID).To(Equal(source.Node.ID))
				})
				It("Should maintain the reference to the node struct", func() {
					opts := &storage.ModelAdapterOpts{
						Source: source,
						Dest:   dest,
					}
					ma := storage.NewModelAdapter(opts)
					err := ma.ExchangeToDest()
					Expect(err).To(BeNil())
					node.ID = 45
					Expect(dest.Node.ID).To(Equal(45))
				})
			})
		})
		Context("Between two models of different types", func() {
			Context("No nested model", func() {
				var source *roach.ChannelConfig
				var dest *storage.ChannelConfig
				BeforeEach(func() {
					source = &roach.ChannelConfig{
						ID:   453,
						Name: "My Channel Config",
					}
					dest = &storage.ChannelConfig{}
				})
				It("Should exchange correctly", func() {
					opts := &storage.ModelAdapterOpts{
						Source:        source,
						Dest:          dest,
						CatalogSource: roach.Catalog(),
						CatalogDest:   storage.Catalog(),
					}
					ma := storage.NewModelAdapter(opts)
					err := ma.ExchangeToDest()
					Expect(err).To(BeNil())
					Expect(source.Node).To(BeNil())
				})
			})
			Context("Nested model", func() {
				var source *roach.ChannelConfig
				var dest *storage.ChannelConfig
				BeforeEach(func() {
					source = &roach.ChannelConfig{
						ID:   453,
						Name: "My Channel Config",
						Node: &roach.Node{
							ID: 96,
						},
					}
					dest = &storage.ChannelConfig{}
				})
				It("Should break ref between old and new nested", func() {
					opts := &storage.ModelAdapterOpts{
						Source:        source,
						Dest:          dest,
						CatalogSource: roach.Catalog(),
						CatalogDest:   storage.Catalog(),
					}
					ma := storage.NewModelAdapter(opts)
					err := ma.ExchangeToDest()
					Expect(err).To(BeNil())
					Expect(source.Node.ID).To(Equal(96))
					Expect(dest.Node.ID).To(Equal(source.Node.ID))
					source.Node.ID = 45
					Expect(dest.Node.ID).To(Equal(96))
				})
			})
			Context("Multiple nested models", func() {
				var source *roach.Range
				var dest *storage.Range
				BeforeEach(func() {
					source = &roach.Range{
						ID: 420,
						LeaseHolderNode: &roach.Node{
							ID: 22,
						},
						ReplicaNodes: []*roach.Node{
							&roach.Node{
								ID: 1,
							},
							&roach.Node{
								ID: 2,
							},
						},
					}
					dest = &storage.Range{
						ID: 1900,
					}
				})
				It("Should exchange correctly", func() {
					opts := &storage.ModelAdapterOpts{
						Source:        source,
						Dest:          dest,
						CatalogSource: roach.Catalog(),
						CatalogDest:   storage.Catalog(),
					}
					ma := storage.NewModelAdapter(opts)
					err := ma.ExchangeToDest()
					Expect(err).To(BeNil())
					Expect(source.ReplicaNodes).To(HaveLen(2))
					Expect(source.ReplicaNodes[0].ID).To(Equal(1))
					Expect(dest.ReplicaNodes[0].ID).To(Equal(source.ReplicaNodes[0].ID))
				})
			})
		})
	})
	Context("Multi model adaptation", func() {
		Context("Models of the same type", func() {
			var source []*roach.ChannelConfig
			var dest []*storage.ChannelConfig
			BeforeEach(func() {
				source = []*roach.ChannelConfig{
					&roach.ChannelConfig{
						ID:   22,
						Name: "Hello",
					},
					&roach.ChannelConfig{
						ID:   24,
						Name: "Hello 24",
					},
				}
				dest = []*storage.ChannelConfig{}
			})
			It("Should exchange correctly", func() {
				opts := &storage.ModelAdapterOpts{
					Source:        &source,
					Dest:          &dest,
					CatalogSource: roach.Catalog(),
					CatalogDest:   storage.Catalog(),
				}
				ma := storage.NewModelAdapter(opts)
				err := ma.ExchangeToDest()
				Expect(err).To(BeNil())
				Expect(dest).To(HaveLen(2))
			})

		})
	})
})
