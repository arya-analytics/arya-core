package model_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Model Exchange", func() {
	Describe("Exchange", func() {
		Context("Single Model Exchange", func() {
			Context("Models of the same type", func() {
				Context("No nested rfl", func() {
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
						me := model.NewExchange(dest, source)
						me.ToSource()
						Expect(source.ID).To(Equal(435))
						Expect(source.ID).To(Equal(dest.ID))
						Expect(source.Name).To(Equal(dest.Name))
					})
					It("Should exchange to dest", func() {
						me := model.NewExchange(source, dest)
						me.ToDest()
						Expect(source.ID).To(Equal(435))
						Expect(source.ID).To(Equal(dest.ID))
						Expect(source.Name).To(Equal(dest.Name))
					})
					It("Shouldn't maintain refs between source and dest models",
						func() {
							me := model.NewExchange(source, dest)
							me.ToDest()
							source.Name = "Hello"
							Expect(dest.Name).To(Equal("Cool Name"))
						})
					It("Should maintain rfl internal refs", func() {
						me := model.NewExchange(source, dest)
						me.ToDest()
						refObj.ID = 9260
						Expect(dest.RefObj.ID).To(Equal(9260))
					})
				})
				Context("With nested rfl", func() {
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
						me := model.NewExchange(dest, source)
						me.ToSource()
						Expect(source.InnerModel.ID).To(Equal(24))
						Expect(dest.InnerModel.ID).To(Equal(source.InnerModel.ID))
					})
					It("Should exchange to dest", func() {
						me := model.NewExchange(source, dest)
						me.ToDest()
						Expect(source.InnerModel.ID).To(Equal(24))
						Expect(dest.InnerModel.ID).To(Equal(source.InnerModel.ID))
					})
					It("Should break the reference to the inner rfl struct", func() {
						me := model.NewExchange(source, dest)
						me.ToDest()
						innerModel.ID = 45
						Expect(dest.InnerModel.ID).To(Equal(24))
					})
				})
			})
			Context("Models of different types", func() {
				Context("No nested rfl", func() {
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
						me := model.NewExchange(source, dest)
						me.ToDest()
						Expect(source.InnerModel).To(BeNil())
					})
				})
				Context("Nested rfl", func() {
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
						me := model.NewExchange(source, dest)
						me.ToDest()
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
						me := model.NewExchange(source, dest)
						me.ToDest()
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
				Context("No nested rfl", func() {
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
						me := model.NewExchange(&source, &dest)
						me.ToDest()
						Expect(dest).To(HaveLen(2))
					})
					It("Should maintain rfl internal refs", func() {
						me := model.NewExchange(&source, &dest)
						me.ToDest()
						Expect(dest[0].RefObj.ID).To(Equal(source[0].RefObj.ID))
						dest[0].RefObj.ID = 720
						Expect(source[0].RefObj.ID).To(Equal(720))
						Expect(source[1].RefObj.ID).To(Equal(720))
					})
				})
				Context("Performing an update", func() {
					BeforeEach(func() {
						source = []*mock.ModelA{
							&mock.ModelA{
								ID:     25,
								Name:   "Hello",
								RefObj: refObj,
							},
							&mock.ModelA{
								ID:     26,
								Name:   "Hello",
								RefObj: refObj,
							},
						}
						dest = []*mock.ModelB{
							&mock.ModelB{
								ID:     25,
								Name:   "Hello 25",
								RefObj: refObj,
							},
							&mock.ModelB{
								ID:     26,
								Name:   "Hello 26",
								RefObj: refObj,
							},
						}
					})
					It("Should set the correct values in source", func() {
						me := model.NewExchange(&source, &dest)
						me.ToSource()
						Expect(dest).To(HaveLen(2))
						Expect(dest[0].Name).To(Equal("Hello 25"))
						Expect(dest[1].Name).To(Equal("Hello 26"))
					})
					Context("Bad Update", func() {
						It("Should warn the caller", func() {
							source[0].ID = 22
							source[1].ID = 28
							me := model.NewExchange(&source, &dest)
							me.ToSource()
							Expect(dest).To(HaveLen(2))
							Expect(dest[0].Name).To(Equal("Hello 25"))
							Expect(dest[1].Name).To(Equal("Hello 26"))
						})
					})
				})
				Context("Multiple nested models", func() {
					var chainInnerModel []*mock.ModelB
					Context("Common chain inner rfl", func() {
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
							me := model.NewExchange(&source, &dest)
							me.ToSource()
							Expect(source).To(HaveLen(2))
							Expect(source[0].CommonChainInnerModel).To(Equal(chainInnerModel))
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
				var me *model.Exchange
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
					me = model.NewExchange(&source, &dest)
					me.ToDest()
				})
				It("Should return the correct source", func() {
					Expect(me.Source().Type()).To(Equal(model.NewReflect(&mock.ModelA{}).
						Type()))
				})
				It("Should return the correct dest", func() {
					Expect(me.Dest().Type()).To(Equal(model.NewReflect(&mock.ModelB{}).Type()))
				})
			})
			Context("Single model", func() {
				var source *mock.ModelA
				var dest *mock.ModelB
				var refObj *mock.RefObj
				var me *model.Exchange
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
					me = model.NewExchange(source, dest)
					me.ToDest()
				})
				It("Should return the correct source", func() {
					Expect(me.Source().Type()).To(Equal(model.NewReflect(&mock.ModelA{}).Type()))
				})
				It("Should return the correct dest", func() {
					Expect(me.Dest().Type()).To(Equal(model.NewReflect(&mock.ModelB{}).Type()))
				})
			})
		})
		Describe("Custom Field Handler", func() {
			Describe("Adapting Primary Keys", func() {
				It("Should adapt a UUID PKC to a string PKC", func() {
					source := &mock.ModelG{
						ID: uuid.New(),
					}
					dest := &mock.ModelH{}
					me := model.NewExchange(source, dest, model.FieldHandlerPK)
					me.ToDest()
					Expect(dest.ID).To(Equal(source.ID.String()))
				})
			})

		})
		Describe("Edge cases + errors", func() {
			Describe("NewStruct Model Adapter", func() {
				Context("Slice and struct mismatch", func() {
					It("Should panic", func() {
						var source []*mock.ModelB
						dest := &mock.ModelB{}
						Expect(func() {
							model.NewExchange(&source, dest)
						}).To(Panic())
					})
				})
				Context("Providing a non-struct or slice type", func() {
					It("Should panic", func() {
						source := &mock.ModelB{}
						dest := 1
						Expect(func() {
							model.NewExchange(source, &dest)
						}).To(Panic())
					})
				})
				Context("Providing a non-pointer Val", func() {
					It("Should panic", func() {
						source := mock.ModelB{}
						dest := &mock.ModelA{}
						Expect(func() {
							model.NewExchange(source, dest)
						}).To(Panic())
					})
				})
				Context("Providing a double-pointer struct Val", func() {
					It("Should panic", func() {
						source := &mock.ModelB{}
						dest := &mock.ModelA{}
						Expect(func() {
							model.NewExchange(&source, &dest)
						}).To(Panic())
					})
				})
				Context("Providing a double-pointer slice Val", func() {
					It("Should panic", func() {
						source := []*mock.ModelB{&mock.ModelB{Name: "Hello"}}
						dest := &[]*mock.ModelA{}
						Expect(func() {
							model.NewExchange(&source, &dest)
						}).To(Panic())
					})
				})
			})
			Describe("Exchanging values", func() {
				Describe("Incompatible rfl types", func() {
					Context("Top level incompatibility", func() {
						Context("Non pointer Val", func() {
							It("Should panic", func() {
								source := &mock.ModelB{
									ID:   22,
									Name: "My Cool Model",
								}
								dest := &mock.ModelC{}
								me := model.NewExchange(source, dest)
								Expect(func() {
									me.ToDest()
								})
							})
						})
						Context("Pointer Val", func() {
							It("Should panic", func() {
								source := &mock.ModelE{
									ID:                  45,
									PointerIncompatible: &map[string]string{"one": "two"},
								}
								dest := &mock.ModelF{}
								me := model.NewExchange(source, dest)
								Expect(func() {
									me.ToDest()
								})
							})
						})
					})
					Context("Nested rfl incompatibility", func() {
						Context("Single rfl", func() {
							Describe("Without the incompatible field defined", func() {
								It("Shouldn't panic", func() {
									source := &mock.ModelD{
										ID: 22,
									}
									dest := &mock.ModelC{}
									me := model.NewExchange(source, dest)
									Expect(func() {
										me.ToDest()
									}).ToNot(Panic())
								})
							})
							Describe("With the incompatible field defined", func() {
								It("Should panic", func() {
									dest := &mock.ModelD{
										ID: 2,
										ModelIncompatible: &mock.ModelB{
											ID:   43,
											Name: "Hello",
										},
									}
									source := &mock.ModelC{}
									me := model.NewExchange(source, dest)
									Expect(func() {
										me.ToSource()
									})
								})
							})
						})
						Context("Chained models", func() {
							It("Should panic", func() {
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
								me := model.NewExchange(source, dest)
								Expect(func() {
									me.ToSource()
								}).To(Panic())
							})
						})
					})
				})
			})

		})
	})

})
