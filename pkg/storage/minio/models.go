package minio

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
)

func catalog() model.Catalog {
	return model.Catalog{
		&channelChunkReplica{},
	}
}

type channelChunkReplica struct {
	ID    uuid.UUID        `model:"role:pk"`
	Telem *telem.ChunkData `model:"role:telemChunkData,"`
}
