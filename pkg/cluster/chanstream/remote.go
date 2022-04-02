package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanstream/v1"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
)

// |||| RPC IMPLEMENTATION ||||

func catalogRemoteRPC() model.Catalog {
	return model.Catalog{
		&api.ChannelSample{},
	}
}

func newExchange(m interface{}) *model.Exchange {
	return rpc.NewModelExchange(m, catalogRemoteRPC().New(m))
}

type RemoteRPC struct {
	rpcPool *cluster.NodeRPCPool
	srp     *remoteStreamRetrievePool
}

func NewRemoteRPC(rpcPool *cluster.NodeRPCPool) *RemoteRPC {
	return &RemoteRPC{
		srp:     newStreamRetrievePool(rpcPool),
		rpcPool: rpcPool,
	}
}

func (r *RemoteRPC) exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&streamq.TSCreate{}:   newStreamCreate(r.rpcPool).exec,
		&streamq.TSRetrieve{}: newRemoteStreamRetrieve(r.srp).exec,
	})
}
