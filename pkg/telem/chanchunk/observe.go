package chanchunk

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"sync"
)

type ObservedChannel struct {
	ConfigPK       uuid.UUID
	LatestChunkPK  uuid.UUID
	ConflictPolicy models.ChannelConflictPolicy
	DataType       telem.DataType
	DataRate       telem.DataRate
	LatestSampleTS telem.TimeStamp
}

type Observe interface {
	//Add(oc ObservedChannel)
	Retrieve(cfgPK uuid.UUID) (ObservedChannel, bool)
	//RetrieveAll() []ObservedChannel
}

type ObserveMem struct {
	mu      sync.Mutex
	chanMap map[uuid.UUID]ObservedChannel
}

func NewObserveMem(channels []ObservedChannel) *ObserveMem {
	chanMap := map[uuid.UUID]ObservedChannel{}
	for _, c := range channels {
		chanMap[c.ConfigPK] = c
	}
	return &ObserveMem{chanMap: chanMap}
}

func (o *ObserveMem) Retrieve(cfgPk uuid.UUID) (ObservedChannel, bool) {
	oc, ok := o.chanMap[cfgPk]
	return oc, ok
}
