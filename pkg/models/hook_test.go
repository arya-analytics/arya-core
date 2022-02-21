package models_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hook", func() {
	Describe("Before Insert", func() {
		Describe("Node GRPC Port Hook", func() {
			It("Should set the GRPC port on the node correctly", func() {
				mRfl := model.NewReflect(&models.Node{
					ID: 1,
				})
				models.HookBeforeNodeInsertSetGRPCPort(mRfl)
				Expect(mRfl.StructFieldByName("GRPCPort").Interface()).To(Equal(models.NodeDefaultGRPCPort))
			})
		})
		It("Should run the correct before insert hooks", func() {
			mRfl := model.NewReflect(&models.Node{})
			Expect(models.BeforeCreate(mRfl)).To(BeNil())
		})
	})
})
