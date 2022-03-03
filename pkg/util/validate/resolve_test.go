package validate_test

import (
	"errors"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func dummyResolveAction(err error, args string) (bool, error) {
	if err.Error() != "weird ol err" {
		return false, err
	}
	if args == "resolveable" {
		return true, nil
	}
	return true, errors.New("unresolveable error")
}

var _ = Describe("Resolve", func() {
	It("should resolve the error successfully", func() {
		run := validate.NewResolve([]func(err error, args string) (bool, error){
			dummyResolveAction,
		})
		err := run.Exec(errors.New("weird ol err"), "resolveable").Error()
		Expect(err).To(BeNil())
		Expect(run.Handled()).To(BeTrue())
		Expect(run.Resolved()).To(BeTrue())
	})
	It("shouldn't resolve the error when its unresolveable", func() {
		run := validate.NewResolve([]func(err error, args string) (bool, error){
			dummyResolveAction,
		})
		err := run.Exec(errors.New("weird ol err"), "unresolveable").Error()
		Expect(err.Error()).To(Equal("unresolveable error"))
		Expect(run.Handled()).To(BeTrue())
		Expect(run.Resolved()).To(BeFalse())
	})
	It("Should return the original error when no resolve can handle it", func() {
		run := validate.NewResolve([]func(err error, args string) (bool, error){
			dummyResolveAction,
		}, errutil.WithAggregation())
		err := run.Exec(errors.New("normal ol err"), "unresolveable").Error()
		Expect(err.Error()).To(Equal("normal ol err"))
		Expect(run.Handled()).To(BeFalse())
		Expect(run.Resolved()).To(BeFalse())
	})

})
