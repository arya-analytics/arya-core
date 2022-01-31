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
	ID        uuid.UUID `cache:"role:tsKey"`
	Name      string
	DataRate  float64
	Retention time.Duration
}

type ChannelSample struct {
	ChannelConfigID uuid.UUID `cache:"role:tsKey,"`
	Value           float64   `cache:"role:tsValue,"`
	Timestamp       int64     `cache:"role:tsStamp"`
}
