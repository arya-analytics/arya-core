package validate

import (
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

type Resolve[T any] struct {
	actions   []func(err error, args T) (bool, error)
	opts      *opts
	catch     errutil.Catch
	sourceErr error
	handled   bool
}

func NewResolve[T any](actions []func(err error, args T) (bool, error), rOpts ...Opt) *Resolve[T] {
	re := &Resolve[T]{
		actions: actions,
		opts:    &opts{},
	}
	for _, opt := range rOpts {
		opt(re.opts)
	}
	return re
}

func (re *Resolve[T]) Exec(err error, args T) *Resolve[T] {
	re.sourceErr = err
	re.catch = &errutil.CatchSimple{}
	if re.opts.aggregate {
		re.catch = &errutil.CatchAggregate{}
	}
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
	return re.Errors()[0]
}
