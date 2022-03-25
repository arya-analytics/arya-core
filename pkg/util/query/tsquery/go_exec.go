package tsquery

import "github.com/arya-analytics/aryacore/pkg/util/query"

type GoExecOpt struct {
	Errors chan error
	Done   chan bool
}

const goExecOptKey query.OptKey = "goExec"

const (
	errorBufferSize = 10
)

func NewGoExecOpt(p *query.Pack) GoExecOpt {
	errors, done := make(chan error, errorBufferSize), make(chan bool)
	o := GoExecOpt{Errors: errors, Done: done}
	p.SetOpt(goExecOptKey, o)
	return o
}

func RetrieveGoExecOpt(p *query.Pack) (GoExecOpt, bool) {
	opt, ok := p.RetrieveOpt(goExecOptKey)
	if !ok {
		return GoExecOpt{}, false
	}
	return opt.(GoExecOpt), true
}
