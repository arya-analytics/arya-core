package chanchunk

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/google/uuid"
	"golang.org/x/sync/semaphore"
	"sync"
)

type observedChannelConfig struct {
	PK     uuid.UUID
	Status models.ChannelStatus
}

type observe interface {
	Add(oc observedChannelConfig)
	Retrieve(pk uuid.UUID) (observedChannelConfig, bool)
}

type observeMem struct {
	sem     *semaphore.Weighted
	mu      sync.Mutex
	chanMap map[uuid.UUID]observedChannelConfig
}

func newObserveMem() *observeMem {
	return &observeMem{chanMap: map[uuid.UUID]observedChannelConfig{}}
}

func (o *observeMem) Retrieve(cfgPk uuid.UUID) (observedChannelConfig, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	oc, ok := o.chanMap[cfgPk]
	return oc, ok
}

func (o *observeMem) Add(oc observedChannelConfig) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.chanMap[oc.PK] = oc
}
