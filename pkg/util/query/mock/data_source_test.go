package mock_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("DataSourceMem", func() {
	var (
		data model.DataSource
		ds   *mock.DataSourceMem
		asm  query.Assemble
		ctx  context.Context
	)
	BeforeEach(func() {
		ctx = context.Background()
		data = model.DataSource{}
		ds = &mock.DataSourceMem{Data: data}
		asm = query.NewAssemble(ds.Exec)
	})

	Describe("Create", func() {
		It("Should add the item to the data", func() {
			Expect(asm.NewCreate().Model(&models.Range{ID: uuid.New()}).Exec(ctx)).To(BeNil())
			Expect(data.Retrieve(reflect.TypeOf(models.Range{})).ChainValue().Interface()).To(HaveLen(1))
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
		//It("Should retrieve by a nested where field", func() {
		//	r := &models.Range{ID: uuid.New()}
		//	rl := &models.RangeLease{RangeID: r.ID, ID: uuid.New()}
		//	Expect(asm.NewCreate().Model(r).Exec(ctx)).To(BeNil())
		//	Expect(asm.NewCreate().Model(rl).Exec(ctx)).To(BeNil())
		//	resR := &models.Range{}
		//	Expect(asm.NewRetrieve().Model(resR).WhereFields(query.WhereFields{"RangeLease.ID": rl.ID}).Exec(ctx)).To(BeNil())
		//	Expect(resR.ID).To(Equal(r.ID))
		//})
		It("Should retrieve the correct relation", func() {
			r := &models.Range{ID: uuid.New()}
			rl := &models.RangeLease{RangeID: r.ID, ID: uuid.New()}
			Expect(asm.NewCreate().Model(r).Exec(ctx)).To(BeNil())
			Expect(asm.NewCreate().Model(rl).Exec(ctx)).To(BeNil())
			resR := &models.Range{}
			Expect(asm.NewRetrieve().Model(resR).WherePK(r.ID).Relation("RangeLease", "ID").Exec(ctx)).To(BeNil())
			Expect(resR.RangeLease.ID).To(Equal(rl.ID))
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
	})
})
