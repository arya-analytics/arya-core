package filter

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

func Exec(p *query.Pack, on interface{}) {
	filtered := p.Model()
	pkc, ok := query.PKOpt(p)
	if ok {
		filtered = execPK(filtered, pkc)
	}
	model.NewReflect(on).Set(filtered)
}

func execPK(sRfl *model.Reflect, pkc model.PKChain) *model.Reflect {
	nRfl := sRfl.NewChain()
	sRfl.ForEach(func(rfl *model.Reflect, i int) {
		for _, pk := range pkc {
			if rfl.PK().Equals(pk) {
				nRfl.ChainAppend(rfl)
			}
		}
	})
	return nRfl
}
