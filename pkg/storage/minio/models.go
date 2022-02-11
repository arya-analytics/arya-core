package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
)

func catalog() model.Catalog {
	return model.Catalog{
		&channelChunkReplica{},
	}
}

type channelChunkReplica struct {
	ID    uuid.UUID `model:"role:pk"`
	Telem storage.Object
}
