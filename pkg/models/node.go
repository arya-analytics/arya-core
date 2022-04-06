package models

import (
	"context"
	"errors"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/uptrace/bun"
	"reflect"
	"strings"
	"time"
)

type Node struct {
	model.Base      `storage:"engines:md,"`
	bun.BaseModel   `bun:"select:nodes_gossip,table:nodes"`
	ID              int       `model:"role:pk" bun:",pk"`
	RPCPort         int       `model:"role:rpc_port" bun:"rpc_port"`
	Address         string    `bun:"type:text,scanonly"`
	IsHost          bool      `bun:"type:boolean,scanonly"`
	StartedAt       time.Time `bun:"type:timestamp,scanonly"`
	IsLive          bool      `bun:"type:boolean,scanonly"`
	Epoch           int       `bun:"type:bigint,scanonly"`
	Expiration      string    `bun:"type:text,scanonly"`
	Draining        bool      `bun:"type:boolean,scanonly"`
	Decommissioning bool      `bun:"type:boolean,scanonly"`
	Membership      string    `bun:"type:text,scanonly"`
	UpdatedAt       time.Time `bun:"type:timestamp,scanonly"`
}

// |||| VALUE ACCESSORS ||||

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

// |||| QUERY HOOK ||||

type NodeQueryHook struct{}

func (nqh *NodeQueryHook) Before(ctx context.Context, p *query.Pack) error {
	qhr := queryHookRunner{rfl: p.Model(), CatchSimple: errutil.NewCatchSimple()}
	if p.Model().Type() == reflect.TypeOf(Node{}) {
		switch p.Query().(type) {
		case *query.Create:
			qhr.Exec(beforeNodeCreateSetDefaultRPCPort)
		}
	}
	return qhr.Error()
}

func (nqh *NodeQueryHook) After(ctx context.Context, p *query.Pack) error {
	return nil
}

const NodeDefaultRPCPort int = 26258

func beforeNodeCreateSetDefaultRPCPort(rfl *model.Reflect) error {
	rfl.ForEach(func(nRfl *model.Reflect, _ int) {
		fld := nRfl.StructFieldByRole(`rpc_port`)
		if fld.IsZero() {
			fld.Set(reflect.ValueOf(NodeDefaultRPCPort))
		}
	})
	return nil
}
