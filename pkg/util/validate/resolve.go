package validate

import (
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

type Resolve[T any] struct {
	actions   []func(err error, args T) (bool, error)
	catch     errutil.Catch
	sourceErr error
	handled   bool
}

func NewResolve[T any](actions []func(err error, args T) (bool, error), opts ...errutil.CatchOpt) *Resolve[T] {
	return &Resolve[T]{actions: actions, catch: errutil.NewCatchSimple(opts...)}
}

func (re *Resolve[T]) Exec(err error, args T) *Resolve[T] {
	re.sourceErr = err
	re.catch.Reset()
	for _, action := range re.actions {
		re.catch.Exec(func() (cErr error) {
			re.handled, cErr = action(err, args)
			return cErr
		})
	}
	return re
}

func (re *Resolve[T]) Handled() bool {
	return re.handled
}

func (re *Resolve[T]) Resolved() bool {
	return re.handled && re.Error() == nil
}

func (re *Resolve[T]) Errors() []error {
	if !re.handled {
		return []error{re.sourceErr}
	}
	return re.catch.Errors()
}

func (re *Resolve[T]) Error() error {
	if len(re.Errors()) == 0 {
		return nil
	}
	return re.Errors()[0]
}
