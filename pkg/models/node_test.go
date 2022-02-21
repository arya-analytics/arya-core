package models_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Node", func() {
	Describe("Value Accessors", func() {
		It("Should return the correct host", func() {
			n := &models.Node{ID: 1, Address: "localhost:26257", RPCPort: models.NodeDefaultRPCPort}
			Expect(n.Host()).To(Equal("localhost"))
		})
		It("Should return the correct RPC address", func() {
			n := &models.Node{ID: 1, Address: "localhost:26257", RPCPort: models.NodeDefaultRPCPort}
			Expect(n.RPCAddress()).To(Equal("localhost:26258"))
		})
	})
	Describe("Query Hook", func() {
		Describe("BeforeQuery", func() {
			Describe("Setting default RPC port", func() {
				It("Should set the default port when none is provided", func() {
					n := &models.Node{ID: 1}
					qh := &models.NodeQueryHook{}
					qe := &storage.QueryEvent{
						Model: model.NewReflect(n),
						Query: &storage.QueryCreate{},
					}
					Expect(qh.BeforeQuery(ctx, qe)).To(BeNil())
					Expect(n.RPCPort).To(Equal(models.NodeDefaultRPCPort))
				})
				It("Shouldn't set the port when a value is provided", func() {
					n := &models.Node{ID: 1, RPCPort: 22}
					qh := &models.NodeQueryHook{}
					qe := &storage.QueryEvent{
						Model: model.NewReflect(n),
						Query: &storage.QueryCreate{},
					}
					Expect(qh.BeforeQuery(ctx, qe)).To(BeNil())
					Expect(n.RPCPort).To(Equal(22))
				})
			})
		})
	})
	Describe("AfterQuery", func() {
		It("Should do nothing", func() {
			n := &models.Node{ID: 1, RPCPort: 22}
			qh := &models.NodeQueryHook{}
			qe := &storage.QueryEvent{
				Model: model.NewReflect(n),
				Query: &storage.QueryCreate{},
			}
			Expect(qh.AfterQuery(ctx, qe)).To(BeNil())
		})
	})

})
