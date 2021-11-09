package live

import (
	"github.com/arya-analytics/aryacore/telem"
	"github.com/google/uuid"
)

type Sender struct {
	id uuid.UUID
	send chan Telemetry
}

type SenderConfig struct {
	channelConfigs []telem.ChannelConfig
}

type TelemetryValue struct {
	timeStamp float64
	value float64
}

type Telemetry map[int32]TelemetryValue