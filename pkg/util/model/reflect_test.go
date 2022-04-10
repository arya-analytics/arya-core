package model_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Reflect", func() {

	Describe("Construction", func() {
		Describe("Pointer Checks", func() {
			It("Should return true when the model is a pointer", func() {
				Expect(model.NewReflect(&mock.ModelA{}).IsPointer()).To(BeTrue())
			})
			It("Should return false when the model is a pointer", func() {
				Expect(model.UnsafeNewReflect(mock.ModelA{}).IsPointer()).To(BeFalse())
			})
		})
		Describe("Pointer Creation", func() {
			It("Should create a new pointer for a non-pointer model", func() {
				Expect(model.UnsafeNewReflect(mock.ModelA{}).ToNewPointer().IsPointer()).To(BeTrue())
			})
			It("Should create the pointer to the correct underlying Val", func() {
				var baseVal []*mock.ModelA
				baseRfl := model.UnsafeNewReflect(baseVal)
				Expect(baseRfl.ToNewPointer().RawValue().Interface()).To(Equal(baseVal))
			})
		})
		Describe("Nesting", func() {
			It("Should avoid wrapping a nested reflect", func() {
				rfl := model.NewReflect(model.NewReflect(&mock.ModelA{}))
				Expect(rfl.IsStruct()).To(BeTrue())
				Expect(rfl.Type()).To(Equal(reflect.TypeOf(mock.ModelA{})))
			})
		})
		Describe("NewReflectFromType", func() {
			It("Should create a reflect from a type", func() {
				rfl := model.NewReflectFromType(reflect.TypeOf(mock.ModelA{}))
				Expect(rfl.IsStruct()).To(BeTrue())
				Expect(rfl.Type()).To(Equal(reflect.TypeOf(mock.ModelA{})))
			})
			It("Should panic if the type is not of struct kind", func() {
				Expect(func() {
					model.NewReflectFromType(reflect.TypeOf(uuid.UUID{}))
				}).To(Panic())
			})
		})
	})

	Context("Single Model", func() {
		var m = &mock.ModelA{
			ID: 22,
			InnerModel: &mock.ModelB{
				ID: 23,
			},
		}
		var mBaseType = reflect.TypeOf(mock.ModelA{})
		var mType = reflect.TypeOf(m)
		var rfl = model.NewReflect(m)
		Describe("The Basics", func() {
			It("Should pass validation without panicking", func() {
				Expect(rfl.Validate).ToNot(Panic())
			})
			It("Should return the correct pointer interface", func() {
				Expect(rfl.Pointer()).To(Equal(m))
			})
			It("Should return the correct type", func() {
				Expect(rfl.Type()).To(Equal(mBaseType))
			})
			It("Should return the correct Val", func() {
				Expect(rfl.StructValue().Type()).To(Equal(mBaseType))
			})
			It("Should return false for IsChain", func() {
				Expect(rfl.IsChain()).To(BeFalse())
			})
			It("Should return true for IsStruct", func() {
				Expect(rfl.IsStruct()).To(BeTrue())
			})
			It("Should panic when trying to access the chan value", func() {
				Expect(func() { rfl.ChanValue() }).To(Panic())
			})
			It("Should set the value without issue", func() {
				nRfl := model.NewReflect(&mock.ModelA{ID: 41})
				Expect(func() {
					nRfl.Set(model.NewReflect(&mock.ModelA{ID: 42}))
				}).ToNot(Panic())
				Expect(nRfl.PK().Raw()).To(Equal(42))

			})
			Describe("Accessing Struct fields", func() {
				It("Should return the correct struct field by name", func() {
					Expect(rfl.StructFieldByName("ID").Interface()).To(Equal(22))
				})
				It("Should be case agnostic when retrieving struct fields", func() {
					Expect(rfl.StructFieldByName("id").Interface()).To(Equal(22))
				})
				It("Should access the correct nested struct field by name", func() {
					Expect(rfl.StructFieldByName("InnerModel.ID").Interface()).To(Equal(23))
				})
				It("Should return a zero reflect when the value does not exist", func() {
					Expect(rfl.StructFieldByName("InnerModel.NonExistent").IsValid()).To(BeFalse())
				})
				It("Should return the correct struct field by role", func() {
					Expect(rfl.StructFieldByRole("pk").Interface()).To(Equal(22))
				})
				It("Should panic if the role doesn't exist", func() {
					Expect(func() {
						rfl.StructFieldByRole("nonexistentrole").Interface()
					}).To(Panic())
				})
				It("Should return the correct struct field by index", func() {
					Expect(rfl.StructValue().Field(0).Interface()).To(Equal(22))
				})
			})

			It("Should return the same item for the raw Val as for the Val",
				func() {
					Expect(rfl.RawValue()).To(Equal(rfl.StructValue()))
				})
			It("Should return the same type for the raw type as for the type", func() {
				Expect(rfl.RawType()).To(Equal(rfl.Type()))
			})
			It("Should return the correct pointer type", func() {
				Expect(rfl.PointerType()).To(Equal(mType))
			})
			It("Should return the correct pointer Val", func() {
				Expect(rfl.PointerValue()).To(Equal(reflect.ValueOf(m)))
			})
			Context("FieldTypeByName", func() {
				It("Should access the correct field type by name", func() {
					Expect(rfl.FieldTypeByName("ID")).To(Equal(reflect.TypeOf(1)))
				})
				It("Should access the correct nested field type by name", func() {
					Expect(rfl.FieldTypeByName("InnerModel.ID")).To(Equal(reflect.TypeOf(1)))
				})
				It("Should panic when the field doesn't exist", func() {
					Expect(func() {
						rfl.FieldTypeByName("InnerModel.NonExistent")
					}).To(PanicWith("field InnerModel.NonExistent does not exist on ModelA"))
				})
			})
		})
		Context("Constructors", func() {
			Describe("New Chain", func() {
				It("Should return the correct type", func() {
					newC := rfl.NewChain()
					Expect(newC.RawType()).To(Equal(reflect.TypeOf([]*mock.ModelA{})))
					Expect(newC.Type()).To(Equal(mBaseType))
				})
			})
			Describe("New Model", func() {
				It("Should return the correct type", func() {
					newM := rfl.NewStruct()
					Expect(newM.PointerType()).To(Equal(mType))
					Expect(newM.Type()).To(Equal(mBaseType))
				})
			})
			Describe("New Raw", func() {
				It("Should return a single model", func() {
					newM := rfl.NewRaw()
					Expect(newM.PointerType()).To(Equal(mType))
					Expect(newM.Type()).To(Equal(mBaseType))
				})
			})
		})
		Describe("For Each", func() {
			It("Should provide the reflect itself", func() {
				rfl.ForEach(func(rfl *model.Reflect, i int) {
					Expect(i).To(Equal(-1))
					Expect(rfl).To(Equal(rfl))
				})
			})
		})
		Describe("PKS", func() {
			It("Should return the correct PKC", func() {
				Expect(rfl.PKChain()).To(HaveLen(1))
				Expect(rfl.PKChain()[0].Raw()).To(Equal(m.ID))
			})
		})
		Describe("Accessing ChainValue", func() {
			It("Should panic", func() {
				Expect(func() {
					rfl.ChainValue()
				}).To(Panic())
			})
		})
	})
	Context("Multiple Models", func() {
		var (
			m               []*mock.ModelA
			mBaseType       reflect.Type
			mType           reflect.Type
			mSingleBaseType reflect.Type
			mSingleType     reflect.Type
			rfl             *model.Reflect
		)
		BeforeEach(func() {
			m = []*mock.ModelA{
				{
					ID: 22,
				},
				{
					ID: 43,
				},
			}
			mBaseType = reflect.TypeOf(m)
			mType = reflect.TypeOf(&m)
			mSingleBaseType = reflect.TypeOf(mock.ModelA{})
			mSingleType = reflect.TypeOf(&mock.ModelA{})
			rfl = model.NewReflect(&m)
		})
		It("Should pass validation without panicking", func() {
			Expect(rfl.Validate).ToNot(Panic())
		})
		It("Should return the correct pointer interface", func() {
			Expect(rfl.Pointer()).To(Equal(&m))
		})
		It("Should return the correct type", func() {
			Expect(rfl.Type()).To(Equal(mSingleBaseType))
		})
		It("Should return true for IsChain", func() {
			Expect(rfl.IsChain()).To(BeTrue())
		})
		It("Should return false for IsStruct", func() {
			Expect(rfl.IsStruct()).To(BeFalse())
		})
		It("Should return the correct model Val by index", func() {
			Expect(rfl.ChainValueByIndex(0).PointerType()).To(Equal(mSingleType))
			Expect(rfl.ChainValueByIndex(0).Type()).To(Equal(mSingleBaseType))
			Expect(rfl.ChainValueByIndexOrNew(0).Pointer()).To(Equal(m[0]))
		})
		It("Should create a new reflect if the index exceeds the chain Val", func() {
			Expect(rfl.ChainValueByIndexOrNew(rfl.ChainValue().Len()).Type()).To(
				Equal(mSingleBaseType))
		})
		It("Should return a slice for the raw Val", func() {
			Expect(rfl.RawValue().Interface()).To(Equal(m))
			Expect(rfl.RawType()).To(Equal(mBaseType))
		})
		It("Should return the correct pointer type", func() {
			Expect(rfl.PointerType()).To(Equal(mType))
		})
		It("Should return the correct pointer Val", func() {
			Expect(rfl.PointerValue()).To(Equal(reflect.ValueOf(&m)))
		})
		It("Should return the correct fields by index", func() {
			Expect(rfl.Fields(0).Raw()[0]).To(Equal(22))
		})
		It("Should return the correct fields by name", func() {
			Expect(rfl.FieldsByName("ID").Raw()[0]).To(Equal(22))
		})
		Describe("Appending to a chain", func() {
			It("Should append to the chain correctly", func() {
				mTwo := &mock.ModelA{ID: 1}
				rfl.ChainAppend(model.NewReflect(mTwo))
				Expect(rfl.ChainValueByIndex(2).Pointer()).To(Equal(mTwo))
			})
			It("Should append each to the chain correctly", func() {
				mTwo := &mock.ModelA{ID: 1}
				rfl.ChainAppendEach(model.NewReflect(mTwo))
				Expect(rfl.ChainValueByIndex(2).Pointer()).To(Equal(mTwo))
			})
		})
		Describe("New Chain", func() {
			It("Should return the correct type", func() {
				newC := rfl.NewChain()
				Expect(newC.RawType()).To(Equal(mBaseType))
				Expect(newC.Type()).To(Equal(mSingleBaseType))
			})
		})
		Describe("New Model", func() {
			It("Should return the correct type", func() {
				newM := rfl.NewStruct()
				Expect(newM.PointerType()).To(Equal(mSingleType))
				Expect(newM.Type()).To(Equal(mSingleBaseType))
			})
		})
		Describe("New Raw", func() {
			It("Should return a model chain", func() {
				newM := rfl.NewRaw()
				Expect(newM.PointerType()).To(Equal(mType))
				Expect(newM.RawType()).To(Equal(mBaseType))
				Expect(newM.Type()).To(Equal(mSingleBaseType))
			})
		})
		Describe("For Each", func() {
			It("Should iterate through the chain", func() {
				rfl.ForEach(func(rflIter *model.Reflect, i int) {
					Expect(rflIter).To(Equal(rfl.ChainValueByIndex(i)))
				})
			})
		})
		Describe("PKS", func() {
			It("Should get the correct Val by PKC", func() {
				val, ok := rfl.ValueByPK(model.NewPK(m[0].ID))
				Expect(ok).To(BeTrue())
				Expect(val.Type()).To(Equal(mSingleBaseType))
			})
			It("Should return not ok when the Val can't be found", func() {
				_, ok := rfl.ValueByPK(model.NewPK(uuid.New()))
				Expect(ok).To(BeFalse())
			})
		})
		Describe("Accessing StructValue", func() {
			It("Should panic", func() {
				Expect(func() {
					rfl.StructValue()
				}).To(Panic())
			})
		})
	})
	Context("Channel of Models", func() {
		var (
			m chan *mock.ModelA
			//mBaseType       reflect.Type
			mSingleBaseType reflect.Type
			//mSingleType     reflect.Type
			rfl *model.Reflect
		)
		BeforeEach(func() {
			m = make(chan *mock.ModelA)
			//mBaseType = reflect.TypeOf(m)
			mSingleBaseType = reflect.TypeOf(mock.ModelA{})
			//mSingleType = reflect.TypeOf(&mock.ModelA{})
			rfl = model.NewReflect(&m)
		})
		It("Should pass validation without panicking", func() {
			Expect(rfl.Validate).ToNot(Panic())
		})
		It("Should return the correct pointer interface", func() {
			Expect(rfl.Pointer()).To(Equal(&m))
		})
		It("Should return the correct type", func() {
			Expect(rfl.Type()).To(Equal(mSingleBaseType))
		})
		It("Should return true for IsChan", func() {
			Expect(rfl.IsChan()).To(BeTrue())
		})
		It("Should return false for IsStruct", func() {
			Expect(rfl.IsStruct()).To(BeFalse())
		})
		It("Should return false for IsChain", func() {
			Expect(rfl.IsChain()).To(BeFalse())
		})
		It("Should panic when trying to get the primary key", func() {
			Expect(func() { rfl.PKChain() }).To(Panic())

		})
		Describe("Sending and receiving", func() {
			It("Should send and receive the values correctly", func() {
				rts := model.NewReflect(&mock.ModelA{ID: 2})
				go rfl.ChanSend(rts)
				r, ok := rfl.ChanRecv()
				Expect(ok).To(BeTrue())
				Expect(r.PK().Raw()).To(Equal(2))
			})
			It("Should send multiple values correctly", func() {
				rts := model.NewReflect(&[]*mock.ModelA{{ID: 1}, {ID: 2}})
				go rfl.ChanSendEach(rts)
				r, ok := rfl.ChanRecv()
				Expect(ok).To(BeTrue())
				Expect(r.PK().Raw()).To(Equal(1))
				r2, ok2 := rfl.ChanRecv()
				Expect(ok2).To(BeTrue())
				Expect(r2.PK().Raw()).To(Equal(2))
			})
		})
	})
	Describe("errors + edge cases", func() {
		It("Should panic when a non pointer is provided", func() {
			Expect(func() {
				model.NewReflect(mock.ModelA{ID: 22})
			}).To(Panic())
		})
		It("Should panic when a non struct is provided", func() {
			i := 11
			Expect(func() {
				model.NewReflect(&i)
			}).To(Panic())
		})
		Context("nil pointer", func() {
			It("Should panic when initializing with a nil struct", func() {
				Expect(func() {
					model.NewReflect((*mock.ModelA)(nil))
				}).To(Panic())
			})
			It("Shouldn't panic when initializing with a nil chain", func() {
				var m []*mock.ModelA
				Expect(func() {
					model.NewReflect(&m)
				}).ToNot(Panic())
			})
		})
	})
	Describe("Filter", func() {
		It("Should filter the reflect into a new reflect", func() {
			id := []*mock.ModelA{{ID: 1}, {ID: 2}}
			od := model.NewReflect(&id).Filter(func(rfl *model.Reflect, i int) bool {
				return rfl.PK().Equals(model.NewPK(2))
			})
			Expect(od.ChainValue().Len()).To(Equal(1))
		})
	})
})
