package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Model Adapter", func() {
	Describe("Adapter", func() {
		Context("Single Model Adaptation", func() {
			Context("Models of the same type", func() {
				Context("No nested refl", func() {
					var source *mock.ModelA
					var dest *mock.ModelA
					var refObj *mock.RefObj
					BeforeEach(func() {
						refObj = &mock.RefObj{
							ID: 220,
						}
						source = &mock.ModelA{
							ID:     435,
							Name:   "Cool Name",
							RefObj: refObj,
						}
						dest = &mock.ModelA{}
					})
					It("Should exchange to source", func() {
						ma := storage.NewModelAdapter(dest, source)
						err := ma.ExchangeToSource()
						Expect(err).To(BeNil())
						Expect(source.ID).To(Equal(435))
						Expect(source.ID).To(Equal(dest.ID))
						Expect(source.Name).To(Equal(dest.Name))
					})
					It("Should exchange to dest", func() {
						ma := storage.NewModelAdapter(source, dest)
						err := ma.ExchangeToDest()
						Expect(err).To(BeNil())
						Expect(source.ID).To(Equal(435))
						Expect(source.ID).To(Equal(dest.ID))
						Expect(source.Name).To(Equal(dest.Name))
					})
					It("Shouldn't maintain refs between source and dest models",
						func() {
							ma := storage.NewModelAdapter(source, dest)
							err := ma.ExchangeToDest()
							if err != nil {
								log.Fatalln(err)
							}
							source.Name = "Hello"
							Expect(dest.Name).To(Equal("Cool Name"))
						})
					It("Should maintain refl internal refs", func() {
						ma := storage.NewModelAdapter(source, dest)
						err := ma.ExchangeToDest()
						if err != nil {
							log.Fatalln(err)
						}
						refObj.ID = 9260
						Expect(dest.RefObj.ID).To(Equal(9260))
					})
				})
				Context("With nested refl", func() {
					var source *mock.ModelA
					var dest *mock.ModelB
					var innerModel *mock.ModelB
					BeforeEach(func() {
						innerModel = &mock.ModelB{
							ID: 24,
						}
						source = &mock.ModelA{
							ID:         420,
							Name:       "Even Cooler Name",
							InnerModel: innerModel,
						}
						dest = &mock.ModelB{}
					})
					It("Should exchange to source", func() {
						ma := storage.NewModelAdapter(dest, source)
						err := ma.ExchangeToSource()
						Expect(err).To(BeNil())
						Expect(source.InnerModel.ID).To(Equal(24))
						Expect(dest.InnerModel.ID).To(Equal(source.InnerModel.ID))
					})
					It("Should exchange to dest", func() {
						ma := storage.NewModelAdapter(source, dest)
						err := ma.ExchangeToDest()
						Expect(err).To(BeNil())
						Expect(source.InnerModel.ID).To(Equal(24))
						Expect(dest.InnerModel.ID).To(Equal(source.InnerModel.ID))
					})
					It("Should break the reference to the inner refl struct", func() {
						ma := storage.NewModelAdapter(source, dest)
						err := ma.ExchangeToDest()
						Expect(err).To(BeNil())
						innerModel.ID = 45
						Expect(dest.InnerModel.ID).To(Equal(24))
					})
				})
			})
			Context("Models of different types", func() {
				Context("No nested refl", func() {
					var source *mock.ModelA
					var dest *mock.ModelB
					BeforeEach(func() {
						source = &mock.ModelA{
							ID:   453,
							Name: "My Channel Config",
						}
						dest = &mock.ModelB{}
					})
					It("Should exchange correctly", func() {
						ma := storage.NewModelAdapter(source, dest)
						err := ma.ExchangeToDest()
						Expect(err).To(BeNil())
						Expect(source.InnerModel).To(BeNil())
					})
				})
				Context("Nested refl", func() {
					var source *mock.ModelB
					var dest *mock.ModelA
					var innerModel *mock.ModelA
					BeforeEach(func() {
						innerModel = &mock.ModelA{
							ID: 96,
						}
						source = &mock.ModelB{
							ID:         453,
							Name:       "My Channel Config",
							InnerModel: innerModel,
						}
						dest = &mock.ModelA{}
					})
					It("Should break ref between old and new nested", func() {
						ma := storage.NewModelAdapter(source, dest)
						err := ma.ExchangeToDest()
						Expect(err).To(BeNil())
						Expect(source.InnerModel.ID).To(Equal(96))
						Expect(dest.InnerModel.ID).To(Equal(source.InnerModel.ID))
						source.InnerModel.ID = 45
						Expect(dest.InnerModel.ID).To(Equal(96))
						innerModel.ID = 32
						Expect(source.InnerModel.ID).To(Equal(32))
					})
				})
				Context("Multiple nested models", func() {
					var source *mock.ModelA
					var dest *mock.ModelB
					BeforeEach(func() {
						source = &mock.ModelA{
							ID: 420,
							InnerModel: &mock.ModelB{
								ID: 22,
							},
							ChainInnerModel: []*mock.ModelB{
								&mock.ModelB{
									ID: 1,
								},
								&mock.ModelB{
									ID: 2,
								},
							},
						}
						dest = &mock.ModelB{
							ID: 1900,
						}
					})
					It("Should exchange correctly", func() {
						ma := storage.NewModelAdapter(source, dest)
						err := ma.ExchangeToDest()
						Expect(err).To(BeNil())
						Expect(source.ID).To(Equal(420))
						Expect(source.ID).To(Equal(dest.ID))
						Expect(source.ChainInnerModel).To(HaveLen(2))
						Expect(source.ChainInnerModel[0].ID).To(Equal(1))
						Expect(dest.ChainInnerModel[0].ID).To(Equal(source.
							ChainInnerModel[0].ID))

					})
				})
			})
		})
		Context("Chain Model Adaptation", func() {
			Context("Models of different types", func() {
				var source []*mock.ModelA
				var dest []*mock.ModelB
				var refObj *mock.RefObj
				Context("No nested refl", func() {
					BeforeEach(func() {
						refObj = &mock.RefObj{
							ID: 672,
						}
						source = []*mock.ModelA{
							&mock.ModelA{
								ID:     22,
								Name:   "Hello",
								RefObj: refObj,
							},
							&mock.ModelA{
								ID:     24,
								Name:   "Hello 24",
								RefObj: refObj,
							},
						}
						dest = []*mock.ModelB{}
					})
					It("Should exchange correctly", func() {
						ma := storage.NewModelAdapter(&source, &dest)
						err := ma.ExchangeToDest()
						Expect(err).To(BeNil())
						Expect(dest).To(HaveLen(2))
					})
					It("Should maintain refl internal refs", func() {
						ma := storage.NewModelAdapter(&source, &dest)
						err := ma.ExchangeToDest()
						if err != nil {
							log.Fatalln(err)
						}
						Expect(dest[0].RefObj.ID).To(Equal(source[0].RefObj.ID))
						dest[0].RefObj.ID = 720
						Expect(source[0].RefObj.ID).To(Equal(720))
						Expect(source[1].RefObj.ID).To(Equal(720))
					})
				})
				Context("Pre-populated dest", func() {
					BeforeEach(func() {
						source = []*mock.ModelA{
							&mock.ModelA{
								ID:     22,
								Name:   "Hello",
								RefObj: refObj,
							},
							&mock.ModelA{
								ID:     24,
								Name:   "Hello 24",
								RefObj: refObj,
							},
						}
						dest = []*mock.ModelB{
							&mock.ModelB{
								ID:     25,
								Name:   "Hello",
								RefObj: refObj,
							},
							&mock.ModelB{
								ID:     26,
								Name:   "Hello 26",
								RefObj: refObj,
							},
						}
					})
					It("Should override the values in dest", func() {
						ma := storage.NewModelAdapter(&source, &dest)
						err := ma.ExchangeToDest()
						if err != nil {
							log.Fatalln(err)
						}
						Expect(dest).To(HaveLen(2))
						Expect(dest[0].ID).To(Equal(22))
						Expect(dest[1].ID).To(Equal(24))
					})
				})
				Context("Multiple nested models", func() {
					var chainInnerModel []*mock.ModelB
					Context("Common chain inner refl", func() {
						BeforeEach(func() {
							refObj = &mock.RefObj{
								ID: 915,
							}
							chainInnerModel = []*mock.ModelB{
								&mock.ModelB{
									ID: 12345,
								},
							}
							source = []*mock.ModelA{}
							dest = []*mock.ModelB{
								&mock.ModelB{
									ID:                    11,
									CommonChainInnerModel: chainInnerModel,
								},
								&mock.ModelB{
									ID:                    12,
									CommonChainInnerModel: chainInnerModel,
								},
							}
						})
						It("Should exchange correctly", func() {
							ma := storage.NewModelAdapter(&source, &dest)
							err := ma.ExchangeToSource()
							Expect(err).To(BeNil())
						})
					})

				})
			})
		})
		Describe("Source and Dest Retrieval", func() {
			Context("Chain of models", func() {
				var source []*mock.ModelA
				var dest []*mock.ModelB
				var refObj *mock.RefObj
				var ma *storage.ModelAdapter
				BeforeEach(func() {
					refObj = &mock.RefObj{
						ID: 672,
					}
					source = []*mock.ModelA{
						&mock.ModelA{
							ID:     22,
							Name:   "Hello",
							RefObj: refObj,
						},
						&mock.ModelA{
							ID:     24,
							Name:   "Hello 24",
							RefObj: refObj,
						},
					}
					dest = []*mock.ModelB{}
					ma = storage.NewModelAdapter(&source, &dest)
					err := ma.ExchangeToDest()
					if err != nil {
						panic(err)
					}
				})
				It("Should return the correct source", func() {
					Expect(ma.Source().Type()).To(Equal(model.NewReflect(&mock.ModelA{}).Type()))
				})
				It("Should return the correct dest", func() {
					Expect(ma.Dest().Type()).To(Equal(model.NewReflect(&mock.ModelB{}).Type()))
				})
			})
			Context("Single model", func() {
				var source *mock.ModelA
				var dest *mock.ModelB
				var refObj *mock.RefObj
				var ma *storage.ModelAdapter
				BeforeEach(func() {
					refObj = &mock.RefObj{
						ID: 672,
					}
					source = &mock.ModelA{
						ID:     22,
						Name:   "Hello",
						RefObj: refObj,
					}
					dest = &mock.ModelB{}
					ma = storage.NewModelAdapter(source, dest)
					err := ma.ExchangeToDest()
					if err != nil {
						panic(err)
					}
				})
				It("Should return the correct source", func() {
					Expect(ma.Source().Type()).To(Equal(model.NewReflect(&mock.ModelA{}).Type()))
				})
				It("Should return the correct dest", func() {
					Expect(ma.Dest().Type()).To(Equal(model.NewReflect(&mock.ModelB{}).Type()))
				})
			})
		})
		Describe("Edge cases + errors", func() {
			Describe("NewModel Model Adapter", func() {
				Context("Slice and struct mismatch", func() {
					It("Should panic", func() {
						var source []*mock.ModelB
						dest := &mock.ModelB{}
						Expect(func() {
							storage.NewModelAdapter(&source, dest)
						}).To(Panic())
					})
				})
				Context("Providing a non-struct or slice type", func() {
					It("Should panic", func() {
						source := &mock.ModelB{}
						dest := 1
						Expect(func() {
							storage.NewModelAdapter(source, &dest)
						}).To(Panic())
					})
				})
				Context("Providing a non-pointer value", func() {
					It("Should panic", func() {
						source := mock.ModelB{}
						dest := &mock.ModelA{}
						Expect(func() {
							storage.NewModelAdapter(source, dest)
						}).To(Panic())
					})
				})
				Context("Providing a double-pointer struct value", func() {
					It("Should panic", func() {
						source := &mock.ModelB{}
						dest := &mock.ModelA{}
						Expect(func() {
							storage.NewModelAdapter(&source, &dest)
						}).To(Panic())
					})
				})
				Context("Providing a double-pointer slice value", func() {
					It("Should return an error", func() {
						source := []*mock.ModelB{&mock.ModelB{Name: "Hello"}}
						dest := &[]*mock.ModelA{}
						Expect(func() {
							storage.NewModelAdapter(&source, &dest)
						}).To(Panic())
					})
				})
			})
			Describe("Exchanging values", func() {
				Describe("Incompatible refl types", func() {
					Context("Top level incompatibility", func() {
						Context("Non pointer value", func() {
							It("Should return an error", func() {
								source := &mock.ModelB{
									ID:   22,
									Name: "My Cool Model",
								}
								dest := &mock.ModelC{}
								ma := storage.NewModelAdapter(source, dest)
								err := ma.ExchangeToDest()
								Expect(err).ToNot(BeNil())
								Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeInvalidField))
							})
						})
						Context("Pointer value", func() {
							It("Should return an error", func() {
								source := &mock.ModelE{
									ID:                  45,
									PointerIncompatible: &map[string]string{"one": "two"},
								}
								dest := &mock.ModelF{}
								ma := storage.NewModelAdapter(source, dest)
								err := ma.ExchangeToDest()
								Expect(err).ToNot(BeNil())
								Expect(err.(storage.Error).Type).To(Equal(storage.
									ErrTypeInvalidField))
							})
						})
					})
					Context("Nested refl incompatibility", func() {
						Context("Single refl", func() {
							Describe("Without the incompatible field defined", func() {
								It("Shouldn't return an error", func() {
									source := &mock.ModelD{
										ID: 22,
									}
									dest := &mock.ModelC{}
									ma := storage.NewModelAdapter(source, dest)
									err := ma.ExchangeToDest()
									Expect(err).To(BeNil())
								})
							})
							Describe("With the incompatible field defined", func() {
								It("Should return an error", func() {
									dest := &mock.ModelD{
										ID: 2,
										ModelIncompatible: &mock.ModelB{
											ID:   43,
											Name: "Hello",
										},
									}
									source := &mock.ModelC{}
									ma := storage.NewModelAdapter(source, dest)
									err := ma.ExchangeToSource()
									Expect(err).ToNot(BeNil())
									Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeInvalidField))
								})
							})
						})
						Context("Chained models", func() {
							It("Should return an error", func() {
								dest := &mock.ModelD{
									ID: 1,
									ChainModelIncompatible: []*mock.ModelB{
										&mock.ModelB{
											ID:   11,
											Name: "String Name",
										},
									},
								}
								source := &mock.ModelC{}
								ma := storage.NewModelAdapter(source, dest)
								err := ma.ExchangeToSource()
								Expect(err).ToNot(BeNil())
								Expect(err.(storage.Error).Type).To(Equal(storage.
									ErrTypeInvalidField))
							})
						})
					})
				})
			})

		})
	})
	Describe("Catalog", func() {
		Context("Valid Catalog", func() {
			Context("Single Model", func() {
				It("Should return a new model of the correct type", func() {
					catalog := storage.ModelCatalog{
						&storage.ChannelConfig{},
					}
					Expect(model.NewReflect(catalog.New(&storage.ChannelConfig{})).
						Pointer()).To(Equal(&storage.ChannelConfig{}))
				})
			})
			Context("Chain of Models", func() {
				It("Should return a new chain of models of the correct type",
					func() {
						catalog := storage.ModelCatalog{
							&storage.ChannelConfig{},
						}
						Expect(model.NewReflect(catalog.New(&[]*storage.ChannelConfig{})).
							Pointer()).To(Equal(&[]*storage.ChannelConfig{}))
					})
			})
		})
		Context("Edge cases + errors", func() {
			It("Should panic when the catalog doesn't contain pointers", func() {
				catalog := storage.ModelCatalog{
					storage.ChannelConfig{},
				}
				Expect(func() {
					catalog.New(storage.ChannelConfig{})
				}).To(Panic())
			})
			It("Should panic when the model cannot be found in the catalog", func() {
				catalog := storage.ModelCatalog{
					&storage.ChannelConfig{},
				}
				Expect(func() {
					catalog.New(&storage.Node{})
				}).To(Panic())
			})
		})
	})
})
