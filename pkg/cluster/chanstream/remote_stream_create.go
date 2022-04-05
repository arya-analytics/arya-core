package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/route"
)

type remoteStreamCreate struct {
	rpcPool *cluster.NodeRPCPool
}

func newRemoteStreamCreate(rpcPool *cluster.NodeRPCPool) *remoteStreamCreate {
	return &remoteStreamCreate{rpcPool: rpcPool}
}

func (s *remoteStreamCreate) newRPC(ctx context.Context, n *models.Node) (api.ChannelStreamService_CreateClient, error) {
	client, err := s.rpcPool.Retrieve(n)
	if err != nil {
		return nil, err
	}
	return api.NewChannelStreamServiceClient(client).Create(ctx)
}

func (s *remoteStreamCreate) exec(ctx context.Context, p *query.Pack) error {
	qStream := stream(p)
	qStream.Segment(func() {
		for {
			rfl, cOk := p.Model().ChanRecv()
			if !cOk || route.CtxDone(ctx) {
				break
			}
			rpcStream, err := s.newRPC(ctx, rfl.StructFieldByName(csFieldNode).Interface().(*models.Node))
			if err != nil {
				qStream.Errors <- err
				break
			}
			exc := newExchange(rfl)
			exc.ToDest()
			if sErr := rpcStream.Send(&api.CreateRequest{Sample: exc.Dest().Pointer().(*api.ChannelSample)}); sErr != nil {
				qStream.Errors <- sErr
			}
		}
	}, streamq.WithSegmentName("cluster.telemstream.remoteStreamCreate"))
	return nil
}
