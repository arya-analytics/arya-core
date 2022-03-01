package cluster

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type StorageService struct {
	storage.Storage
}

func NewStorageService(store storage.Storage) *StorageService {
	return &StorageService{Storage: store}
}

func (ss *StorageService) CanHandle(q *query.Pack) bool {
	return model.Catalog{
		&models.Node{},
		&models.Range{},
		&models.RangeReplica{},
		&models.RangeLease{},
		&models.ChannelConfig{},
		&models.ChannelChunk{},
	}.Contains(q.Model().Pointer())
}
