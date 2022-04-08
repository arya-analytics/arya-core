package query_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Memo", func() {
	var (
		pkc  model.PKChain
		memo *query.Memo
		exec = &mock.Exec{}
	)
	BeforeEach(func() {
		memo = query.NewMemo(model.NewReflect(&[]*models.ChannelConfig{}))
		for i := 0; i < 5; i++ {
			id := uuid.New()
			pkc = append(pkc, model.NewPK(id))
			memo.Add(model.NewReflect(&models.ChannelConfig{ID: id, Name: "Hello"}))
		}
	})
	It("Should return the results of the memoized query", func() {
		var resCC []*models.ChannelConfig
		Expect(query.
			NewRetrieve().
			Model(&resCC).
			WithMemo(memo).
			BindExec(exec.Exec).
			WherePKs(pkc).
			Exec(context.Background())).
			To(BeNil())
		Expect(resCC).To(HaveLen(5))
	})
	It("Shouldn't return the results of the memoized query when no pk is provided", func() {
		var resCC []*models.ChannelConfig
		Expect(query.
			NewRetrieve().
			Model(&resCC).
			WithMemo(memo).
			BindExec(exec.Exec).
			Exec(context.Background())).
			To(BeNil())
		Expect(resCC).To(HaveLen(0))
	})
	It("Should add the results of a query to the memo", func() {
		newID := uuid.New()
		resCC := &models.ChannelConfig{}
		Expect(query.
			NewRetrieve().
			Model(resCC).
			WherePK(newID).
			WithMemo(memo).
			BindExec(func(ctx context.Context, p *query.Pack) error {
				p.Model().Set(model.NewReflect(&models.ChannelConfig{ID: newID, Name: "Hello"}))
				return nil
			}).
			Exec(ctx),
		)
		Expect(resCC.Name).To(Equal("Hello"))
		resCCTwo := &models.ChannelConfig{}
		Expect(query.
			NewRetrieve().
			Model(resCCTwo).
			WherePK(newID).
			WithMemo(memo).
			BindExec(exec.Exec).
			Exec(ctx),
		)
		Expect(resCCTwo.Name).To(Equal("Hello"))
	})
})
