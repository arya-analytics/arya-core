package validate_test

import (
	"errors"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type DummyResolve struct{}

func (dr *DummyResolve) CanHandle(err error) bool {
	return err.Error() == "weird ol err"
}

func (dr *DummyResolve) Handle(err error, args interface{}) error {
	if args == "resolveable" {
		return nil
	}
	return errors.New("unresolveable error")
}

var _ = Describe("Resolve", func() {
	It("should resolve the error successfully", func() {
		resolves := []validate.Resolve{
			&DummyResolve{},
		}
		run := validate.NewResolveRun(resolves)
		err := run.Exec(errors.New("weird ol err"), "resolveable").Error()
		Expect(err).To(BeNil())
		Expect(run.Handled()).To(BeTrue())
		Expect(run.Resolved()).To(BeTrue())
	})
	It("shouldn't resolve the error when its unresolveable", func() {
		resolves := []validate.Resolve{
			&DummyResolve{},
		}
		run := validate.NewResolveRun(resolves)
		err := run.Exec(errors.New("weird ol err"), "unresolveable").Error()
		Expect(err.Error()).To(Equal("unresolveable error"))
		Expect(run.Handled()).To(BeTrue())
		Expect(run.Resolved()).To(BeFalse())
	})
	It("Should return the original error when no resolve can handle it", func() {
		resolves := []validate.Resolve{
			&DummyResolve{},
		}
		run := validate.NewResolveRun(resolves, validate.WithAggregation())
		err := run.Exec(errors.New("normal ol err"), "unresolveable").Error()
		Expect(err.Error()).To(Equal("normal ol err"))
		Expect(run.Handled()).To(BeFalse())
		Expect(run.Resolved()).To(BeFalse())
	})

})
