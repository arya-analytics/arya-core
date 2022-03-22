package query_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	modelmock "github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Query", func() {
	Describe("Model Binding", func() {
		It("Should allow the caller to bind a reflect", func() {
			m := &modelmock.ModelA{}
			p := query.NewRetrieve().Model(model.NewReflect(m)).Pack()
			Expect(p.Model().Pointer()).To(Equal(m))
		})
	})
	Describe("Switch", func() {
		DescribeTable("Executing the correct query", func(q query.Query, expected int) {
			actual := 0
			Expect(query.Switch(ctx, q.Pack(), query.Ops{
				&query.Create{}: func(ctx context.Context, p *query.Pack) error {
					actual = 1
					return nil
				},
				&query.Retrieve{}: func(ctx context.Context, p *query.Pack) error {
					actual = 2
					return nil
				},
				&query.Update{}: func(ctx context.Context, p *query.Pack) error {
					actual = 3
					return nil
				},
				&query.Delete{}: func(ctx context.Context, p *query.Pack) error {
					actual = 4
					return nil
				},
				&query.Migrate{}: func(ctx context.Context, p *query.Pack) error {
					actual = 5
					return nil
				},
			})).To(BeNil())
			Expect(actual).To(Equal(expected))
		},
			Entry("Create Query", query.NewCreate(), 1),
			Entry("Retrieve Query", query.NewRetrieve(), 2),
			Entry("Update Query", query.NewUpdate(), 3),
			Entry("Delete Query", query.NewDelete(), 4),
			Entry("Migrate Query", query.NewMigrate(), 5),
		)
		It("Should panic when there is no viable query handler", func() {
			Expect(func() {
				query.Switch(ctx, query.NewRetrieve().Pack(), query.Ops{})
			}).To(Panic())
		})
	})
	Describe("Pack", func() {
		Describe("String", func() {
			It("Should return the query as a string", func() {
				p := query.NewRetrieve().Model(&[]*models.ChannelConfig{}).Pack()
				Expect(p.String()).
					To(Equal("Variant: *query.Retrieve, Model: models.ChannelConfig, Count: 0, Opts: map[]"))
			})
		})
	})

})
