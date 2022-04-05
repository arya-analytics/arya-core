package route_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Batch", func() {
	Describe("BatchModel", func() {
		It("Should batch a set of models correctly", func() {
			models := []interface{}{
				&mock.ModelA{
					ID:   1,
					Name: "test1",
				},
				&mock.ModelA{
					ID:   2,
					Name: "test3",
				},
				&mock.ModelA{
					ID:   32,
					Name: "test3",
				},
			}

			b := route.BatchModel[string](model.NewReflect(&models), "Name")
			Expect(b["test3"].ChainValue().Len()).To(Equal(2))
			Expect(b["test1"].ChainValue().Len()).To(Equal(1))
		})
		Describe("Edge Cases + Errors", func() {
			It("Should panic when a field has the wrong type", func() {
				models := []interface{}{
					&mock.ModelA{
						ID:   1,
						Name: "test1",
					},
					&mock.ModelA{
						ID:   2,
						Name: "test3",
					},
					&mock.ModelA{
						ID:   32,
						Name: "test3",
					},
				}
				Expect(func() {
					route.BatchModel[int](model.NewReflect(&models), "Name")
				}).To(Panic())
			})
		})
	})
	Describe("ModelSwitchIter", func() {
		It("Should iterate an action over each group of models", func() {
			models := []*mock.ModelA{
				{
					ID:   1,
					Name: "test1",
				},
				{
					ID:   2,
					Name: "test3",
				},
				{
					ID:   32,
					Name: "test3",
				},
			}
			calledTimes := 0
			route.ModelSwitchIter[string](
				model.NewReflect(&models),
				"Name",
				func(fld string, m *model.Reflect) {
					calledTimes++
				},
			)
			Expect(calledTimes).To(Equal(2))
		})
	})
	Describe("ModelSwitchBoolean", func() {
		It("Should switch over a boolean field", func() {
			models := []*mock.ModelA{
				{
					ID:           1,
					Name:         "test1",
					BooleanField: true,
				},
				{
					ID:           2,
					Name:         "test3",
					BooleanField: true,
				},
				{
					ID:           32,
					Name:         "test3",
					BooleanField: false,
				},
			}
			trueCalledTimes := 0
			falseCalledTimes := 0
			route.ModelSwitchBoolean(
				model.NewReflect(&models),
				"BooleanField",
				func(m *model.Reflect) {
					trueCalledTimes++
				},
				func(rfl *model.Reflect) {
					falseCalledTimes++
				},
			)
			Expect(falseCalledTimes).To(Equal(1))
			Expect(trueCalledTimes).To(Equal(1))
		})
	})
})
