package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
	"strings"
	"time"
)

type Node struct {
	ID              int `model:"role:pk"`
	Address         string
	RPCPort         int `model:"role:rpc_port"`
	StartedAt       time.Time
	IsLive          bool
	IsHost          bool
	Epoch           int
	Expiration      string
	Draining        bool
	Decommissioning bool
	Membership      string
	UpdatedAt       time.Time
}

func (n *Node) Host() string {
	sn := strings.Split(n.Address, ":")
	return sn[0]
}

func (n *Node) RPCAddress() (string, error) {
	if n.Host() == "" || n.RPCPort == 0 {
		return "", errors.New("node provided no address or grpc port")
	}
	return fmt.Sprintf("%s:%v", n.Host(), n.RPCPort), nil
}

type NodeQueryHook struct{}

func (nqh *NodeQueryHook) BeforeQuery(ctx context.Context, qe *storage.QueryEvent) error {
	qhr := queryHookRunner{rfl: qe.Model, Catcher: &errutil.Catcher{}}
	if qe.Model.Type() == reflect.TypeOf(Node{}) {
		switch qe.Query.(type) {
		case *storage.QueryCreate:
			qhr.Exec(beforeNodeCreateSetDefaultRPCPort)
		}
	}
	return qhr.Error()
}

func (nqh *NodeQueryHook) AfterQuery(ctx context.Context, qe *storage.QueryEvent) error {
	return nil
}

const NodeDefaultRPCPort = 26258

func beforeNodeCreateSetDefaultRPCPort(rfl *model.Reflect) error {
	rfl.ForEach(func(nRfl *model.Reflect, _ int) {
		fld := nRfl.StructFieldByRole(`rpc_port`)
		if fld.IsZero() {
			fld.Set(reflect.ValueOf(NodeDefaultRPCPort))
		}
	})
	return nil
}
