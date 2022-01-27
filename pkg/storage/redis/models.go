package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	"time"
)

// |||| CATALOG ||||
var _catalog = storage.ModelCatalog{
	&ChannelConfig{},
	&ChannelSample{},
}

func catalog() storage.ModelCatalog {
	return _catalog
}

type ChannelConfig struct {
	ID        uuid.UUID
	Name      string
	DataRate  float64
	Retention time.Duration
}

type ChannelSample struct {
	Value           float32
	Timestamp       time.Time
	ChannelConfigID uuid.UUID
}
