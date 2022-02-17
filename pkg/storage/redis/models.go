package redis

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"time"
)

// |||| CATALOG ||||

func catalog() model.Catalog {
	return model.Catalog{
		&channelConfig{},
		&channelSample{},
	}
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
