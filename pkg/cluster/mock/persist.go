package mock

import (
	"github.com/arya-analytics/aryacore/pkg/util/query"
	querymock "github.com/arya-analytics/aryacore/pkg/util/query/mock"
)

type Persist struct {
	*querymock.DataSourceMem
}

func (pst *Persist) CanHandle(p *query.Pack) bool {
	return true
}
