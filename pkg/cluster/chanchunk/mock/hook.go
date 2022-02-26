package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

type HostInterceptQueryHook int

func (h HostInterceptQueryHook) AfterQuery(ctx context.Context, qe *storage.QueryEvent) error {
	switch qe.Model.Type() {
	case reflect.TypeOf(models.ChannelChunkReplica{}):
		h.setNodeIsHost(qe.Model, "RangeReplica.Node.ID", "RangeReplica.Node.IsHost")
	case reflect.TypeOf(models.RangeReplica{}):
		h.setNodeIsHost(qe.Model, "Node.ID", "Node.IsHost")
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

func (h HostInterceptQueryHook) BeforeQuery(ctx context.Context, qe *storage.QueryEvent) error {
	return nil
}
