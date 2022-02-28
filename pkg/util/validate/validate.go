package validate

import "github.com/arya-analytics/aryacore/pkg/util/errutil"

// |||| VALIDATE ||||

type Validate[T any] struct {
	actions []func(T) error
	opts    *opts
	catch   errutil.Catch
}

func New[T any](actions []func(T) error, vOpts ...Opt) *Validate[T] {
	v := &Validate[T]{actions: actions, opts: &opts{}}
	for _, opt := range vOpts {
		opt(v.opts)
	}
	return v
}

func (v *Validate[T]) Exec(m T) *Validate[T] {
	v.catch = &errutil.CatchSimple{}
	if v.opts.aggregate {
		v.catch = &errutil.CatchAggregate{}
	}
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

// |||| OPTIONS |||

type opts struct {
	aggregate bool
}

type Opt func(vo *opts)

func WithAggregation() Opt {
	return func(vo *opts) {
		vo.aggregate = true
	}
}
