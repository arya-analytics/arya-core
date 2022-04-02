package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
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
	srp     *streamRetrievePool
}

func NewRemoteRPC(rpcPool *cluster.NodeRPCPool) *RemoteRPC {
	return &RemoteRPC{
		srp:     newStreamRetrievePool(context.Background(), rpcPool),
		rpcPool: rpcPool,
	}
}

func (r *RemoteRPC) exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(ctx, p, query.Ops{
		&streamq.TSCreate{}:   r.create,
		&streamq.TSRetrieve{}: newStreamRetrieve(r.srp).exec,
	})
}

func (r *RemoteRPC) create(ctx context.Context, p *query.Pack) error {
	s := stream(p)
	s.Segment(func() {
		for {
			rfl, cOk := p.Model().ChanRecv()
			if !cOk || route.CtxDone(ctx) {
				break
			}
			stream, err := r.newCreateStream(ctx, rfl.StructFieldByName(csFieldNode).Interface().(*models.Node))
			if err != nil {
				s.Errors <- err
				break
			}
			exc := newExchange(rfl)
			exc.ToDest()
			if sErr := stream.Send(&api.CreateRequest{Sample: exc.Dest().Pointer().(*api.ChannelSample)}); sErr != nil {
				s.Errors <- sErr
			}
		}
	})
	return nil
}

func (r *RemoteRPC) newCreateStream(ctx context.Context, n *models.Node) (api.ChannelStreamService_CreateClient, error) {
	client, err := r.rpcPool.Retrieve(n)
	if err != nil {
		return nil, err
	}
	return api.NewChannelStreamServiceClient(client).Create(ctx)
}
