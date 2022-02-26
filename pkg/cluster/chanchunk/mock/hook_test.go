package mock_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hook", func() {
	It("Should intercept the query and set IsHost to false", func() {
		store.AddQueryHook(mock.HostInterceptQueryHook(2))
		Expect(store.NewCreate().Model(&models.Node{ID: 1}).Exec(ctx)).To(BeNil())
		rng := &models.Range{}
		Expect(store.NewCreate().Model(rng).Exec(ctx)).To(BeNil())
		RR := &models.RangeReplica{
			RangeID: rng.ID,
			NodeID:  1,
		}
		Expect(store.NewCreate().Model(RR).Exec(ctx)).To(BeNil())
		resRR := &models.RangeReplica{}
		Expect(store.NewRetrieve().Model(resRR).Relation("Node", "ID", "IsHost").WherePK(RR.ID).Exec(ctx))
		Expect(resRR.Node.IsHost).To(BeFalse())
	})
})
