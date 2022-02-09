package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
)

func catalog() model.Catalog {
	return model.Catalog{
		&channelChunk{},
	}
}

type channelChunk struct {
	ID   uuid.UUID `model:"role:pk"`
	Data storage.Object
}
