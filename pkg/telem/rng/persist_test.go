package rng_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Persist", func() {
	Describe("PersistCluster", func() {
		var (
			newRng *models.Range
			p      rng.Persist
			items  []interface{}
		)
		BeforeEach(func() {
			if clust == nil {
				var err error
				clust, err = mock.New(ctx)
				Expect(err).To(BeNil())
			}
			p = &rng.PersistCluster{Cluster: clust}
			node := &models.Node{ID: 1}
			newRng = &models.Range{ID: uuid.New()}
			items = []interface{}{
				node,
				newRng,
			}
		})
		JustBeforeEach(func() {
			for _, item := range items {
				Expect(clust.NewCreate().Model(item).Exec(ctx)).To(BeNil())
			}
		})
		JustAfterEach(func() {
			for _, item := range items {
				Expect(clust.NewDelete().Model(item).WherePKs(model.NewReflect(item).PKChain().Raw()).Exec(ctx)).To(BeNil())
			}
		})
		Describe("NewRange", func() {
			It("Should save a new range, range lease, and range replica to storage", func() {
				rng, err := p.NewRange(ctx, 1)
				Expect(err).To(BeNil())
				Expect(model.NewPK(rng.ID).IsZero()).To(BeFalse())
				Expect(model.NewPK(rng.RangeLease.ID).IsZero()).To(BeFalse())
				Expect(model.NewPK(rng.RangeLease.RangeReplica.ID).IsZero()).To(BeFalse())
			})
		})
		Describe("New Range Replica", func() {
			It("Should save the replica with the correct node id", func() {
				p := &rng.PersistCluster{Cluster: clust}
				rngReplica, err := p.NewRangeReplica(ctx, newRng.ID, 1)
				Expect(err).To(BeNil())
				Expect(rngReplica.NodeID).To(Equal(1))
			})
		})
		Describe("Retrieve Range Replica", func() {
			It("Should retrieve the correct replica", func() {
				rngReplica, err := p.NewRangeReplica(ctx, newRng.ID, 1)
				Expect(err).To(BeNil())
				rr, err := p.RetrieveRangeReplicas(ctx, newRng.ID)
				Expect(err).To(BeNil())
				Expect(rr).To(HaveLen(1))
				Expect(rr[0].ID).To(Equal(rngReplica.ID))
			})
		})
		Describe("Retrieve a Range", func() {
			It("Should retrieve the correct range", func() {
				resRng, err := p.RetrieveRange(ctx, newRng.ID)
				Expect(err).To(BeNil())
				Expect(resRng.ID).To(Equal(newRng.ID))
			})
		})
		Describe("Retrieve Range Chunk Replicas", func() {

		})
	})

})
