package cluster_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryDelete", func() {
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
		It("Should bind the correct model", func() {
			Expect(clus.NewDelete().Model(m).Exec(ctx))
			Expect(svc.QueryRequest.Model.Pointer().(*storage.ChannelChunkReplica).ID).To(Equal(m.ID))
		})
		Context("WherePK", func() {
			It("Should bind the correct PK", func() {
				pk := uuid.New()
				Expect(clus.NewDelete().Model(m).WherePK(pk).Exec(ctx))
				pkOpt, ok := cluster.PKQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(pkOpt).To(Equal(model.NewPKChain([]uuid.UUID{pk})))
			})
		})
		Context("WherePKs", func() {
			It("Should bind the correct PKs", func() {
				pks := model.NewPKChain([]uuid.UUID{uuid.New(), uuid.New()})
				Expect(clus.NewDelete().Model(m).WherePKs(pks.Raw()...).Exec(ctx))
				pkOpt, ok := cluster.PKQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(pkOpt).To(Equal(pks))
			})
		})
	})
})
