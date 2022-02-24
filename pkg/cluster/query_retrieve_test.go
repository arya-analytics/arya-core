package cluster_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
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
				Expect(clus.NewRetrieve().Model(m).WherePK(pk).Exec(ctx)).To(BeNil())
				pkOpt, ok := internal.PKQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(pkOpt).To(Equal(model.NewPKChain([]uuid.UUID{pk})))
			})
		})
		Context("WherePKs", func() {
			It("Should bind the correct PKs", func() {
				pks := model.NewPKChain([]uuid.UUID{uuid.New(), uuid.New()})
				Expect(clus.NewRetrieve().Model(m).WherePKs(pks.Raw()).Exec(ctx)).To(BeNil())
				pkOpt, ok := internal.PKQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(pkOpt).To(Equal(pks))
			})
		})
		Context("WhereFields", func() {
			It("Should set the correct fields", func() {
				flds := model.WhereFields{"key": "value"}
				Expect(clus.NewRetrieve().Model(m).WhereFields(flds).Exec(ctx)).To(BeNil())
				fldOpt, ok := internal.WhereFieldsQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(fldOpt).To(Equal(flds))
			})
		})
		Context("Relation", func() {
			It("Should set the correct relations", func() {
				Expect(clus.NewRetrieve().Model(m).Relation("rel", "fldOne").Exec(ctx)).To(BeNil())
				relOpts := internal.RelationQueryOpts(svc.QueryRequest)
				Expect(relOpts).To(HaveLen(1))
				Expect(relOpts[0].Rel).To(Equal("rel"))
				Expect(relOpts[0].Fields).To(Equal([]string{"fldOne"}))
			})
			It("Should allow for multiple relations", func() {
				Expect(clus.NewRetrieve().
					Model(m).
					Relation("rel", "fldOne").
					Relation("relTwo", "fldTwo").
					Exec(ctx),
				).To(BeNil())
				relOpts := internal.RelationQueryOpts(svc.QueryRequest)
				Expect(relOpts).To(HaveLen(2))
				Expect(relOpts[0].Rel).To(Equal("rel"))
				Expect(relOpts[0].Fields).To(Equal([]string{"fldOne"}))
				Expect(relOpts[1].Rel).To(Equal("relTwo"))
				Expect(relOpts[1].Fields).To(Equal([]string{"fldTwo"}))
			})
		})
		Context("Fields", func() {
			It("Should set the correct fields", func() {
				Expect(clus.NewRetrieve().Model(m).Fields("ID", "RandomField").Exec(ctx)).To(BeNil())
				fldOpt, ok := internal.RetrieveFieldsQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(fldOpt).To(Equal(internal.FieldsQueryOpt{"ID", "RandomField"}))

			})
		})
		Context("Calculations", func() {
			It("Should set the correct calculations", func() {
				Expect(clus.NewRetrieve().Model(m).Calculate(storage.CalculateAVG, "ID", 0).Exec(ctx)).To(BeNil())
				calcOpt, ok := internal.RetrieveCalculateQueryOpt(svc.QueryRequest)
				Expect(ok).To(BeTrue())
				Expect(calcOpt.Into).To(Equal(0))
				Expect(calcOpt.FldName).To(Equal("ID"))
				Expect(calcOpt.C).To(Equal(storage.CalculateAVG))
			})
		})
	})
})
