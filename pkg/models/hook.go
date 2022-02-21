package models

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

type Hook func(rfl *model.Reflect) error

// || NODE RELATED ||

const NodeDefaultGRPCPort = 26258

func HookBeforeNodeInsertSetGRPCPort(rfl *model.Reflect) error {
	if rfl.Type() == reflect.TypeOf(Node{}) {
		rfl.ForEach(func(nRfl *model.Reflect, _ int) {
			fld := nRfl.StructFieldByRole(`grpc_port`)
			if fld.IsZero() {
				fld.Set(reflect.ValueOf(NodeDefaultGRPCPort))
			}
		})
	}
	return nil
}

func BeforeCreate(rfl *model.Reflect) error {
	hooks := []Hook{
		HookBeforeNodeInsertSetGRPCPort,
	}
	for _, h := range hooks {
		if err := h(rfl); err != nil {
			return err
		}
	}
	return nil
}
