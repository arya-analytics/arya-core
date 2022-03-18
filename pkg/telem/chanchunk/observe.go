package chanchunk

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/google/uuid"
	"golang.org/x/sync/semaphore"
	"sync"
)

type ObservedChannelConfig struct {
	PK     uuid.UUID
	Status models.ChannelStatus
}

type Observe interface {
	Add(oc ObservedChannelConfig)
	Retrieve(pk uuid.UUID) (ObservedChannelConfig, bool)
}

type ObserveMem struct {
	sem     *semaphore.Weighted
	mu      sync.Mutex
	chanMap map[uuid.UUID]ObservedChannelConfig
}

func NewObserveMem() *ObserveMem {
	return &ObserveMem{chanMap: map[uuid.UUID]ObservedChannelConfig{}}
}

func (o *ObserveMem) Retrieve(cfgPk uuid.UUID) (ObservedChannelConfig, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	oc, ok := o.chanMap[cfgPk]
	return oc, ok
}

func (o *ObserveMem) Add(oc ObservedChannelConfig) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.chanMap[oc.PK] = oc
}
