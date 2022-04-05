package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"reflect"
)

type HostInterceptQueryHook int

func (h HostInterceptQueryHook) After(ctx context.Context, p *query.Pack) error {
	switch p.Query().(type) {
	case *query.Retrieve:
		switch p.Model().Type() {
		case reflect.TypeOf(models.ChannelChunkReplica{}):
			h.setNodeIsHost(p.Model(), "RangeReplica.Node.ID", "RangeReplica.Node.IsHost")
		case reflect.TypeOf(models.RangeReplica{}):
			h.setNodeIsHost(p.Model(), "Node.ID", "Node.IsHost")
		case reflect.TypeOf(models.ChannelConfig{}):
			h.setNodeIsHost(p.Model(), "Node.ID", "Node.IsHost")
		}
	}
	return nil
}

func (h HostInterceptQueryHook) setNodeIsHost(rfl *model.Reflect, nodePKFld, nodeIsHostFld string) {
	rfl.ForEach(func(nRfl *model.Reflect, i int) {
		pk := nRfl.StructFieldByName(nodePKFld)
		if !pk.IsValid() {
			return
		}
		pko := pk.Interface()
		nih := nRfl.StructFieldByName(nodeIsHostFld)
		if pko == int(h) {
			nih.Set(reflect.ValueOf(true))
		} else {
			nih.Set(reflect.ValueOf(false))
		}
	})
}

func (h HostInterceptQueryHook) Before(ctx context.Context, p *query.Pack) error {
	return nil
}
