package validate

import "github.com/arya-analytics/aryacore/pkg/util/errutil"

// |||| VALIDATE ||||

type Validate[T any] struct {
	actions []func(T) error
	catch   *errutil.CatchSimple
}

func New[T any](actions []func(T) error, opts ...errutil.CatchOpt) *Validate[T] {
	v := &Validate[T]{actions: actions, catch: errutil.NewCatchSimple(opts...)}
	return v
}

func (v *Validate[T]) Exec(m T) *Validate[T] {
	v.catch.Reset()
	for _, action := range v.actions {
		v.catch.Exec(func() error { return action(m) })
	}
	return v
}

func (v *Validate[T]) Error() error {
	return v.catch.Error()
}

func (v *Validate[T]) Errors() []error {
	return v.catch.Errors()
}
