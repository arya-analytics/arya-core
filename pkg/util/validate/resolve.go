package validate

import (
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
)

type Resolve interface {
	CanHandle(err error) bool
	Handle(err error, args interface{}) error
}

type ResolveRun struct {
	resolves  []Resolve
	opts      *opts
	catch     errutil.Catch
	sourceErr error
	handled   bool
}

func NewResolveRun(resolves []Resolve, rOpts ...Opt) *ResolveRun {
	re := &ResolveRun{
		resolves: resolves,
		opts:     &opts{},
	}
	for _, opt := range rOpts {
		opt(re.opts)
	}
	return re
}

func (re *ResolveRun) Exec(err error, args interface{}) *ResolveRun {
	re.sourceErr = err
	re.catch = &errutil.CatchSimple{}
	if re.opts.aggregate {
		re.catch = &errutil.CatchAggregate{}
	}
	for _, resolve := range re.resolves {
		if resolve.CanHandle(err) {
			re.handled = true
			re.catch.Exec(func() error { return resolve.Handle(err, args) })
		}
	}
	return re
}

func (re *ResolveRun) Handled() bool {
	return re.handled
}

func (re *ResolveRun) Resolved() bool {
	return re.handled && re.Error() == nil
}

func (re *ResolveRun) Errors() []error {
	if !re.handled {
		return []error{re.sourceErr}
	}
	return re.catch.Errors()
}

func (re *ResolveRun) Error() error {
	return re.Errors()[0]
}
