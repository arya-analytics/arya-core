package mock_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("DataSourceMem", func() {
	var (
		ds  *mock.DataSourceMem
		asm query.Assemble
		ctx context.Context
	)
	BeforeEach(func() {
		ctx = context.Background()
		ds = mock.NewDataSourceMem()
		asm = query.NewAssemble(ds.Exec)
	})

	Describe("Create", func() {
		It("Should add the item to the data", func() {
			Expect(asm.NewCreate().Model(&models.Range{ID: uuid.New()}).Exec(ctx)).To(BeNil())
			Expect(ds.Data.Retrieve(reflect.TypeOf(models.Range{})).ChainValue().Interface()).To(HaveLen(1))
		})
	})
	Describe("Retrieve", func() {
		It("Should retrieve by the item's PK", func() {
			cR := &models.Range{ID: uuid.New()}
			Expect(asm.NewCreate().Model(cR).Exec(ctx)).To(BeNil())
			resR := &models.Range{}
			Expect(asm.NewRetrieve().Model(resR).WherePK(cR.ID).Exec(ctx)).To(BeNil())
			Expect(resR.ID).To(Equal(cR.ID))
		})
		It("Should retrieve by the item's where field", func() {
			cR := []*models.Range{{ID: uuid.New()}, {ID: uuid.New()}}
			Expect(asm.NewCreate().Model(&cR).Exec(ctx)).To(BeNil())
			var resR []*models.Range
			Expect(asm.NewRetrieve().Model(&resR).WhereFields(query.WhereFields{"ID": cR[0].ID}).Exec(ctx)).To(BeNil())
			Expect(resR).To(HaveLen(1))
			Expect(resR[0].ID).To(Equal(cR[0].ID))
		})
		It("Should retrieve by a nested where field", func() {
			r := &models.Range{ID: uuid.New()}
			rl := &models.RangeLease{RangeID: r.ID, ID: uuid.New()}
			Expect(asm.NewCreate().Model(r).Exec(ctx)).To(BeNil())
			Expect(asm.NewCreate().Model(&models.Range{ID: uuid.New()}).Exec(ctx)).To(BeNil())
			Expect(asm.NewCreate().Model(rl).Exec(ctx)).To(BeNil())
			var resR []*models.Range
			Expect(asm.NewRetrieve().Model(&resR).WhereFields(query.WhereFields{"RangeLease.ID": rl.ID}).Exec(ctx)).To(BeNil())
			Expect(resR).To(HaveLen(1))
			Expect(resR[0].ID).To(Equal(r.ID))
		})
		It("Should retrieve the correct relation", func() {
			r := &models.Range{ID: uuid.New()}
			rl := &models.RangeLease{RangeID: r.ID, ID: uuid.New()}
			Expect(asm.NewCreate().Model(r).Exec(ctx)).To(BeNil())
			Expect(asm.NewCreate().Model(rl).Exec(ctx)).To(BeNil())
			resR := &models.Range{}
			Expect(asm.NewRetrieve().Model(resR).WherePK(r.ID).Relation("RangeLease", "ID").Exec(ctx)).To(BeNil())
			Expect(resR.RangeLease.ID).To(Equal(rl.ID))
		})
		It("Should panic when a relation can't be found on an item", func() {
			r := &models.Range{ID: uuid.New()}
			Expect(asm.NewCreate().Model(r).Exec(ctx)).To(BeNil())
			Expect(func() {
				asm.NewRetrieve().Model(&models.Range{}).WherePK(r.ID).Relation("IDontExist", "ID").Exec(ctx)
			}).To(Panic())
		})
		It("Should retrieve the correct nested relation", func() {
			r := &models.Range{ID: uuid.New()}
			rr := &models.RangeReplica{NodeID: 1}
			rl := &models.RangeLease{RangeID: r.ID, ID: uuid.New(), RangeReplicaID: rr.ID}
			Expect(asm.NewCreate().Model(r).Exec(ctx)).To(BeNil())
			Expect(asm.NewCreate().Model(rr).Exec(ctx)).To(BeNil())
			Expect(asm.NewCreate().Model(rl).Exec(ctx)).To(BeNil())
			resR := &models.Range{}
			Expect(asm.NewRetrieve().Model(resR).WherePK(r.ID).Relation("RangeLease.RangeReplica", "ID").Exec(ctx)).To(BeNil())
			Expect(resR.RangeLease.RangeReplica.ID).To(Equal(rr.ID))
			Expect(resR.RangeLease.RangeReplica.NodeID).To(Equal(1))
		})
		It("Should calculate a value correctly", func() {
			cc := []*models.ChannelChunk{
				{
					ID:   uuid.New(),
					Size: 10000,
				},
				{
					ID:   uuid.New(),
					Size: 10000,
				},
			}
			Expect(asm.NewCreate().Model(&cc).Exec(ctx)).To(BeNil())
			var size int64
			Expect(asm.NewRetrieve().Model(&models.ChannelChunk{}).Calc(query.CalcSum, "Size", &size).Exec(ctx))
			Expect(size).To(Equal(int64(20000)))
		})
	})
	Describe("Update", func() {
		Describe("Unary", func() {
			It("Should update a single item correctly", func() {
				r := &models.Range{ID: uuid.New()}
				rr := &models.RangeReplica{ID: uuid.New(), NodeID: 1}
				rl := &models.RangeLease{ID: uuid.New(), RangeID: r.ID, RangeReplicaID: rr.ID}
				Expect(asm.NewCreate().Model(r).Exec(ctx)).To(BeNil())
				Expect(asm.NewCreate().Model(rr).Exec(ctx)).To(BeNil())
				Expect(asm.NewCreate().Model(rl).Exec(ctx)).To(BeNil())
				updateRR := &models.RangeReplica{NodeID: 2}
				Expect(asm.NewUpdate().Model(updateRR).WherePK(rr.ID).Fields("NodeID").Exec(ctx)).To(BeNil())
				resRR := &models.RangeReplica{}
				Expect(asm.NewRetrieve().Model(resRR).WherePK(rr.ID).Exec(ctx)).To(BeNil())
				Expect(resRR.NodeID).To(Equal(2))
			})
			It("Should return an error when the item isn't found", func() {
				err := asm.NewUpdate().Model(&models.RangeReplica{}).WherePK(uuid.New()).Fields("NodeID").Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeItemNotFound))
			})
			It("Should panic when a field isn't found", func() {
				r := &models.Range{ID: uuid.New()}
				rr := &models.RangeReplica{ID: uuid.New(), NodeID: 1}
				rl := &models.RangeLease{ID: uuid.New(), RangeID: r.ID, RangeReplicaID: rr.ID}
				Expect(asm.NewCreate().Model(r).Exec(ctx)).To(BeNil())
				Expect(asm.NewCreate().Model(rr).Exec(ctx)).To(BeNil())
				Expect(asm.NewCreate().Model(rl).Exec(ctx)).To(BeNil())
				updateRR := &models.RangeReplica{NodeID: 2}
				Expect(func() {
					_ = asm.NewUpdate().Model(updateRR).WherePK(rr.ID).Fields("IDontExist").Exec(ctx)
				}).To(Panic())
			})
			It("Should panic when no fields are specified", func() {
				r := &models.Range{ID: uuid.New()}
				rr := &models.RangeReplica{ID: uuid.New(), NodeID: 1}
				rl := &models.RangeLease{ID: uuid.New(), RangeID: r.ID, RangeReplicaID: rr.ID}
				Expect(asm.NewCreate().Model(r).Exec(ctx)).To(BeNil())
				Expect(asm.NewCreate().Model(rr).Exec(ctx)).To(BeNil())
				Expect(asm.NewCreate().Model(rl).Exec(ctx)).To(BeNil())
				updateRR := &models.RangeReplica{NodeID: 2}
				Expect(func() {
					_ = asm.NewUpdate().Model(updateRR).WherePK(rr.ID).Exec(ctx)
				}).To(Panic())
			})
		})
		Describe("Bulk", func() {
			It("Should bulk update items correctly", func() {
				r := &models.Range{ID: uuid.New()}
				rr := []*models.RangeReplica{
					{
						ID:     uuid.New(),
						NodeID: 1,
					},
					{
						ID:     uuid.New(),
						NodeID: 2,
					},
				}
				Expect(asm.NewCreate().Model(r).Exec(ctx)).To(BeNil())
				Expect(asm.NewCreate().Model(&rr).Exec(ctx)).To(BeNil())
				updateRR := []*models.RangeReplica{
					{
						ID:     rr[0].ID,
						NodeID: 3,
					},
					{
						ID:     rr[1].ID,
						NodeID: 4,
					},
				}
				Expect(asm.NewUpdate().Model(&updateRR).Fields("NodeID").Bulk().Exec(ctx)).To(BeNil())
				resRR := &models.RangeReplica{}
				Expect(asm.NewRetrieve().Model(resRR).WherePK(rr[0].ID).Exec(ctx)).To(BeNil())
				Expect(resRR.NodeID).To(Equal(3))
				resRRTwo := &models.RangeReplica{}
				err := asm.NewRetrieve().Model(resRRTwo).WhereFields(query.WhereFields{"NodeID": 2}).Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeItemNotFound))
			})
		})
	})
	Describe("Delete", func() {
		It("Should delete an item correctly", func() {
			n := &models.Node{ID: 1}
			Expect(asm.NewCreate().Model(n).Exec(ctx)).To(BeNil())
			Expect(asm.NewDelete().Model(n).WherePK(n.ID).Exec(ctx)).To(BeNil())
			err := asm.NewRetrieve().Model(n).WherePK(n.ID).Exec(ctx)
			Expect(err).ToNot(BeNil())
			Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeItemNotFound))
		})
		It("Should return a not found error when the item to delete is not found", func() {
			err := asm.NewDelete().Model(&models.ChannelSample{}).WherePK(uuid.New()).Exec(ctx)
			Expect(err).ToNot(BeNil())
			Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeItemNotFound))
		})
	})
})
