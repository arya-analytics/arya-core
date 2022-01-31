package redis

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	"time"
)

// |||| CATALOG ||||
var _catalog = storage.ModelCatalog{
	&channelConfig{},
	&channelSample{},
}

func catalog() storage.ModelCatalog {
	return _catalog
}

type channelConfig struct {
	ID        uuid.UUID `model:"role:tsKey"`
	Name      string
	DataRate  float64
	Retention time.Duration
}

type channelSample struct {
	ChannelConfigID uuid.UUID `model:"role:tsKey,"`
	Value           float64   `model:"role:tsValue,"`
	Timestamp       int64     `model:"role:tsStamp"`
}
