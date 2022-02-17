package cluster_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryCreate", func() {
	Context("Query Assembly", func() {
		var (
			m    *storage.ChannelChunkReplica
			clus cluster.Cluster
			svc  *mock.Service
			ctx  = context.Background()
		)
		BeforeEach(func() {
			svc = mock.NewService()
			clus = cluster.New(cluster.ServiceChain{
				svc,
			})
			m = &storage.ChannelChunkReplica{
				ID: uuid.New(),
			}
		})
		JustBeforeEach(func() {
			err := clus.NewCreate().Model(m).Exec(ctx)
			Expect(err).To(BeNil())
		})
		It("Should bind the correct model", func() {
			Expect(svc.QueryRequest.Model.Pointer().(*storage.ChannelChunkReplica).ID).To(Equal(m.ID))
		})
		It("Should have the correct query variant", func() {
			Expect(svc.QueryRequest.Variant).To(Equal(internal.QueryVariantCreate))
		})
	})
})
