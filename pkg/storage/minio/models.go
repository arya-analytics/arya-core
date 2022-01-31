package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
)

var _catalog = storage.ModelCatalog{
	&channelChunk{},
}

func catalog() storage.ModelCatalog {
	return _catalog
}

type channelChunk struct {
	ID   uuid.UUID `model:"role:pk"`
	Data storage.Object
}
