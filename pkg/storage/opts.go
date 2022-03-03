package storage

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type OptConverters []OptConverter

type OptConverter func(p *query.Pack)

func (ocs OptConverters) Exec(p *query.Pack) {
	for _, oc := range ocs {
		oc(p)
	}
}
