package rng_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Persist", func() {
	Describe("PersistCluster", func() {
		var node *models.Node
		BeforeEach(func() {
			if clust == nil {
				var err error
				clust, err = mock.New(ctx)
				Expect(err).To(BeNil())
			}
		})
		Describe("NewRange", func() {
			BeforeEach(func() {
				node = &models.Node{ID: 1}
			})
			JustBeforeEach(func() {
				Expect(clust.NewCreate().Model(node).Exec(ctx)).To(BeNil())
			})
			JustAfterEach(func() {
				Expect(clust.NewDelete().Model(node).WherePK(node.ID).Exec(ctx)).To(BeNil())
			})
			It("Should save a new range, range lease, and range replica to storage", func() {
				p := &rng.PersistCluster{Cluster: clust}
				rng, err := p.NewRange(ctx, 1)
				Expect(err).To(BeNil())
				Expect(model.NewPK(rng.ID).IsZero()).To(BeFalse())
				Expect(model.NewPK(rng.RangeLease.ID).IsZero()).To(BeFalse())
				Expect(model.NewPK(rng.RangeLease.RangeReplica.ID).IsZero()).To(BeFalse())
			})
		})
	})

})
