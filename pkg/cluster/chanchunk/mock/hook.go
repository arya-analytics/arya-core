package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"reflect"
)

type HostInterceptQueryHook int

func (h HostInterceptQueryHook) AfterQuery(ctx context.Context, p *query.Pack) error {
	switch p.Model().Type() {
	case reflect.TypeOf(models.ChannelChunkReplica{}):
		h.setNodeIsHost(p.Model(), "RangeReplica.Node.ID", "RangeReplica.Node.IsHost")
	case reflect.TypeOf(models.RangeReplica{}):
		h.setNodeIsHost(p.Model(), "Node.ID", "Node.IsHost")
	}
	return nil
}

func (h HostInterceptQueryHook) setNodeIsHost(rfl *model.Reflect, nodePKFld, nodeIsHostFld string) {
	rfl.ForEach(func(nRfl *model.Reflect, i int) {
		nodePKFld := nRfl.StructFieldByName(nodePKFld)
		if !nodePKFld.IsValid() {
			return
		}
		nodePK := nodePKFld.Interface()
		nodeIsHost := nRfl.StructFieldByName(nodeIsHostFld)
		if nodePK == int(h) {
			nodeIsHost.Set(reflect.ValueOf(true))
		} else {
			nodeIsHost.Set(reflect.ValueOf(false))
		}
	})
}

func (h HostInterceptQueryHook) BeforeQuery(ctx context.Context, p *query.Pack) error {
	return nil
}
