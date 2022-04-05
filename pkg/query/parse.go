package query

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
)

func ParsePKC(strPKC []string) (model.PKChain, error) {
	pkc := model.NewPKChain([]uuid.UUID{})
	for _, strPK := range strPKC {
		pk, err := model.NewPK(uuid.New()).NewFromString(strPK)
		if err != nil {
			return pkc, err
		}
		pkc = append(pkc, pk)
	}
	return pkc, nil
}
