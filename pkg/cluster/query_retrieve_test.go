package cluster_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/arya-analytics/aryacore/pkg/cluster"
)

var _ = Describe("QueryRetrieve", func() {
	Context("Query Assembly", func() {
		var (
			m    *models.ChannelChunkReplica
			clus cluster.Cluster
			svc  *mock.Service
			ctx  = context.Background()
		)
		BeforeEach(func() {
			svc = mock.NewService()
			clus = cluster.New(cluster.ServiceChain{
				svc,
			})
			m = &models.ChannelChunkReplica{
				ID: uuid.New(),
			}
		})
		It("Should bind the correct model", func() {
			Expect(clus.NewRetrieve().Model(m).Exec(ctx))
			Expect(svc.QueryRequest.Model.Pointer().(*models.ChannelChunkReplica).ID).To(Equal(m.ID))
		})
		Context("WherePK", func() {
			It("Should bind the correct PK", func() {
				pk := uuid.New()
				Expect(clus.NewRetrieve().Model(m).WherePK(pk).Exec(ctx))
				pkOpt, ok := internal.PKQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(pkOpt).To(Equal(model.NewPKChain([]uuid.UUID{pk})))
			})
		})
		Context("WherePKs", func() {
			It("Should bind the correct PKs", func() {
				pks := model.NewPKChain([]uuid.UUID{uuid.New(), uuid.New()})
				Expect(clus.NewRetrieve().Model(m).WherePKs(pks.Raw()).Exec(ctx))
				pkOpt, ok := internal.PKQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(pkOpt).To(Equal(pks))
			})
		})
		Context("WhereFields", func() {
			It("Should set the correct fields", func() {
				flds := models.Fields{"key": "value"}
				Expect(clus.NewRetrieve().Model(m).WhereFields(flds).Exec(ctx))
				fldOpt, ok := internal.FieldsQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(fldOpt).To(Equal(flds))
			})
		})
	})
})
