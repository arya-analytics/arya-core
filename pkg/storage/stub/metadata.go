package stub

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
)

type MDEngine struct{}

func (e *MDEngine) NewAdapter() storage.Adapter {
	return &mdAdapter{id: uuid.New()}
}

func (e *MDEngine) IsAdapter(a storage.Adapter) bool {
	_, ok := e.bindAdapter(a)
	return ok
}

func (e *MDEngine) bindAdapter(a storage.Adapter) (*mdAdapter, bool) {
	ma, ok := a.(*mdAdapter)
	return ma, ok
}

type mdAdapter struct {
	id uuid.UUID
}

func (a *mdAdapter) ID() uuid.UUID {
	return a.id
}
