package models_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
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
		Describe("Before", func() {
			Describe("Setting default RPC port", func() {
				It("Should set the default port when none is provided", func() {
					n := &models.Node{ID: 1}
					qh := &models.NodeQueryHook{}
					p := query.NewCreate().Model(n).Pack()
					Expect(qh.Before(ctx, p)).To(BeNil())
					Expect(n.RPCPort).To(Equal(models.NodeDefaultRPCPort))
				})
				It("Shouldn't set the port when a value is provided", func() {
					n := &models.Node{ID: 1, RPCPort: 22}
					qh := &models.NodeQueryHook{}
					p := query.NewCreate().Model(n).Pack()
					Expect(qh.Before(ctx, p)).To(BeNil())
					Expect(n.RPCPort).To(Equal(22))
				})
			})
		})
	})
	Describe("After", func() {
		It("Should do nothing", func() {
			n := &models.Node{ID: 1, RPCPort: 22}
			qh := &models.NodeQueryHook{}
			p := query.NewRetrieve().Model(n).Pack()
			Expect(qh.After(ctx, p)).To(BeNil())
		})
	})

})
