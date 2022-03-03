package redis

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
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
	DataRate  telem.DataRate
	Retention time.Duration
}

type channelSample struct {
	ChannelConfigID uuid.UUID       `model:"role:tsKey,"`
	Value           float64         `model:"role:tsValue,"`
	Timestamp       telem.TimeStamp `model:"role:tsStamp"`
}
