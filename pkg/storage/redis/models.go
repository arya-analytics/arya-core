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
	ID        uuid.UUID `model:"role:tsKey"`
	Name      string
	DataRate  float64
	Retention time.Duration
}

type ChannelSample struct {
	ChannelConfigID uuid.UUID `model:"role:tsKey,"`
	Value           float64   `model:"role:tsValue,"`
	Timestamp       int64     `model:"role:tsStamp"`
}
