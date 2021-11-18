package live

import (
	"fmt"
	"github.com/arya-analytics/aryacore/ds"
)

type Locator interface {
	Locate(ChanCfgIds []int32)
}

func NewLocator(pooler *ds.ConnPooler) Locator {
	return &DatabaseLocator{pooler}
}

type DatabaseLocator struct {
	pooler *ds.ConnPooler
}

func (dl DatabaseLocator) Locate(ChanCfgIds []int32) {
	fmt.Println("Locating")
}
