package base_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/base"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	var (
		svc *base.Service
	)
	BeforeEach(func() {
		svc = base.NewService(store)
	})
	Describe("Can Handle", func() {
		It("Should return true for all of the following models", func() {
			validModels := []interface{}{
				&models.Node{},
				&models.Range{},
				&models.RangeReplica{},
				&models.RangeLease{},
				&models.ChannelConfig{},
			}
			for _, m := range validModels {
				qr := internal.NewQueryRequest(internal.QueryVariantRetrieve, model.NewReflect(m))
				Expect(svc.CanHandle(qr)).To(BeTrue())
			}
		})
	})
	Describe("Create + Delete + Retrieve + Update Queries", func() {
		Context("Standard usage", func() {
			var (
				node          *models.Node
				channelConfig *models.ChannelConfig
			)
			BeforeEach(func() {
				node = &models.Node{ID: 1}
				channelConfig = &models.ChannelConfig{NodeID: node.ID}
			})
			AfterEach(func() {

			})
			It("Should create, update, retrieve, and delete the items correctly", func() {
				// Creation
				By("Creating the node")
				nodeRfl := model.NewReflect(node)
				ccRfl := model.NewReflect(channelConfig)
				nodeCreateQR := internal.NewQueryRequest(internal.QueryVariantCreate, nodeRfl)
				Expect(svc.Exec(ctx, nodeCreateQR)).To(BeNil())
				By("Creating the channel config")
				ccCreateQR := internal.NewQueryRequest(internal.QueryVariantCreate, ccRfl)
				Expect(svc.Exec(ctx, ccCreateQR)).To(BeNil())

				// Update
				channelConfig.Name = "Cool Name"
				ccUpdateQR := internal.NewQueryRequest(internal.QueryVariantUpdate, ccRfl)
				internal.NewPKQueryOpt(ccUpdateQR, channelConfig.ID)
				Expect(svc.Exec(ctx, ccUpdateQR)).To(BeNil())

				// Retrieve
				By("Retrieving the channel config by PK")
				ccRetrieveByPKRes := &models.ChannelConfig{}
				ccRetrieveByPKResRfl := model.NewReflect(ccRetrieveByPKRes)
				ccRetrieveByPKQR := internal.NewQueryRequest(internal.QueryVariantRetrieve, ccRetrieveByPKResRfl)
				internal.NewPKQueryOpt(ccRetrieveByPKQR, channelConfig.ID)
				internal.NewFieldsQueryOpt(ccRetrieveByPKQR, "ID", "Name", "NodeID", "DataRate")
				internal.NewRelationQueryOpt(ccRetrieveByPKQR, "Node", "ID")
				Expect(svc.Exec(ctx, ccRetrieveByPKQR)).To(BeNil())
				Expect(ccRetrieveByPKResRfl.PK().Raw()).To(Equal(channelConfig.ID))
				Expect(ccRetrieveByPKRes.Node.ID).To(Equal(node.ID))

				By("Retrieving the channel config by the node PK")
				ccRetrieveByNodeIDResRfl := model.NewReflect(&models.ChannelConfig{})
				ccRetrieveByNodeIDQR := internal.NewQueryRequest(internal.QueryVariantRetrieve, ccRetrieveByNodeIDResRfl)
				internal.NewWhereFieldsQueryOpt(ccRetrieveByNodeIDQR, model.WhereFields{"Node.ID": 1})
				Expect(svc.Exec(ctx, ccRetrieveByNodeIDQR)).To(BeNil())
				Expect(ccRetrieveByNodeIDResRfl.PK().Raw()).To(Equal(channelConfig.ID))

				By("Deleting the channel config")
				ccDeleteQR := internal.NewQueryRequest(internal.QueryVariantDelete, ccRfl)
				internal.NewPKQueryOpt(ccDeleteQR, channelConfig.ID)
				Expect(svc.Exec(ctx, ccDeleteQR)).To(BeNil())
				By("Deleting the node")
				nodeDeleteQR := internal.NewQueryRequest(internal.QueryVariantDelete, nodeRfl)
				internal.NewPKQueryOpt(nodeDeleteQR, node.ID)
				Expect(svc.Exec(ctx, nodeDeleteQR)).To(BeNil())
			})
		})
	})
	Describe("Retrieve Query", func() {
		Context("By primary key", func() {})
	})

})
