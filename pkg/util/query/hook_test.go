package query_test

import (
	"context"
	"errors"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type MockQueryHook struct{}

func (m *MockQueryHook) Before(ctx context.Context, p *query.Pack) error {
	return nil
}

func (m *MockQueryHook) After(ctx context.Context, p *query.Pack) error {
	_, ok := query.RetrievePKOpt(p)
	if ok {
		return errors.New("A really badass error")
	}
	return nil
}

var _ = Describe("Hook", func() {
	var (
		h  = &MockQueryHook{}
		hr = query.NewHookRunner()
	)
	Describe("HookRunner", func() {
		Describe("AddQueryHook", func() {
			It("Should add the query hook correctly", func() {
				Expect(func() { hr.AddQueryHook(h) }).ToNot(Panic())
			})
		})
		Describe("ClearQueryHooks", func() {
			It("Should clear the query hooks correctly", func() {
				hr.AddQueryHook(h)
				Expect(func() {
					hr.ClearQueryHooks()
				}).ToNot(Panic())
			})
		})
		Describe("RemoveQueryHook", func() {
			It("Should remove the query hooks correctly", func() {
				hr.AddQueryHook(h)
				Expect(func() {
					hr.RemoveQueryHook(h)
				}).ToNot(Panic())
			})
		})
		Describe("Before", func() {
			It("Should run the before hooks without error", func() {
				hr.AddQueryHook(h)
				Expect(hr.Before(context.Background(), &query.Pack{})).To(Succeed())
			})
		})
		Describe("After", func() {
			It("Should run the after hooks without error", func() {
				hr.AddQueryHook(h)
				Expect(hr.After(context.Background(), &query.Pack{})).To(Succeed())
			})
		})
		Describe("Should catch errors", func() {
			It("Should catch errors encountered while running after hooks", func() {
				hr.AddQueryHook(h)
				Expect(hr.After(context.Background(), query.NewRetrieve().WherePK(12).Pack())).ToNot(Succeed())
			})
		})
	})
})
