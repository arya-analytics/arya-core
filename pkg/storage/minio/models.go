package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
)

func catalog() storage.ModelCatalog {
	return storage.ModelCatalog{
		&channelChunk{},
	}
}

type channelChunk struct {
	ID   uuid.UUID `model:"role:pk"`
	Data storage.Object
}
