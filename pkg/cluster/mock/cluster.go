package mock

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
)

func NewCluster() *cluster.Cluster {
	return cluster.New(
		mock.NewStorage(),
	)
}
