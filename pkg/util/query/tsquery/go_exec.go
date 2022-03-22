package tsquery

import "github.com/arya-analytics/aryacore/pkg/util/query"

type GoExecOpt struct {
	Errors chan error
}

const goExecOptKey query.OptKey = "goExec"

func NewGoExecOpt(p *query.Pack, e chan error) {
	p.SetOpt(goExecOptKey, GoExecOpt{Errors: e})
}

func RetrieveGoExecOpt(p *query.Pack) (GoExecOpt, bool) {
	opt, ok := p.RetrieveOpt(goExecOptKey)
	if !ok {
		return GoExecOpt{}, false
	}
	return opt.(GoExecOpt), true
}
