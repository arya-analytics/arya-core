package chanchunk

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/google/uuid"
	"sync"
)

type ObservedChannelConfig struct {
	PK    uuid.UUID
	State models.ChannelState
}

type Observe interface {
	Add(oc ObservedChannelConfig)
	Retrieve(pk uuid.UUID) (ObservedChannelConfig, bool)
}

type ObserveMem struct {
	mu      sync.Mutex
	chanMap map[uuid.UUID]ObservedChannelConfig
}

func NewObserveMem() *ObserveMem {
	return &ObserveMem{chanMap: map[uuid.UUID]ObservedChannelConfig{}}
}

func (o *ObserveMem) Retrieve(cfgPk uuid.UUID) (ObservedChannelConfig, bool) {
	oc, ok := o.chanMap[cfgPk]
	return oc, ok
}

func (o *ObserveMem) Add(oc ObservedChannelConfig) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.chanMap[oc.PK] = oc
}
