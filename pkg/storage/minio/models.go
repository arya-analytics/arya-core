package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
)

var _catalog = storage.ModelCatalog{
	&ChannelChunk{},
}

func catalog() storage.ModelCatalog {
	return _catalog
}

type ChannelChunk struct {
	ID   uuid.UUID `model:"role:pk"`
	Data storage.Object
}
