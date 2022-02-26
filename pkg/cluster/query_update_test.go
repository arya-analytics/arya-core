package cluster_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryUpdate", func() {
	Context("Query Assembly", func() {
		var (
			m    *models.ChannelChunkReplica
			clus cluster.Cluster
			svc  *mock.Service
			ctx  = context.Background()
		)
		BeforeEach(func() {
			svc = mock.NewService()
			clus = cluster.New()
			clus.BindService(svc)
			m = &models.ChannelChunkReplica{
				ID: uuid.New(),
			}
		})
		It("Should bind the correct model", func() {
			Expect(clus.NewUpdate().Model(m).Exec(ctx))
			Expect(svc.QueryRequest.Model.Pointer().(*models.ChannelChunkReplica).ID).To(Equal(m.ID))
		})
		Context("WherePK", func() {
			It("Should bind the correct PK", func() {
				pk := uuid.New()
				Expect(clus.NewUpdate().Model(m).WherePK(pk).Exec(ctx))
				pkOpt, ok := internal.PKQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(pkOpt).To(Equal(model.NewPKChain([]uuid.UUID{pk})))
			})
		})
		Context("Bulk", func() {
			It("Should set the correct bulk option", func() {
				Expect(clus.NewUpdate().Model(m).Bulk().Exec(ctx))
				bulkOpt := internal.BulkUpdateQueryOpt(svc.QueryRequest)
				Expect(bulkOpt).To(BeTrue())
			})
		})
		Context("Fields", func() {
			It("Should set the correct fields", func() {
				Expect(clus.NewUpdate().Model(m).Fields("Name").Exec(ctx))
				fieldsOpt, ok := internal.RetrieveFieldsQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(fieldsOpt).To(Equal(internal.FieldsQueryOpt{"Name"}))
			})
		})
	})
})
