package bulk

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
)

type Bulk struct {
	cluster cluster.Cluster
}

func (b *Bulk) CreateStream() *Create {

}
