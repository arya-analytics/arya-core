package chanstream

import (
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"io"
)

type ServerRPC struct {
	api.UnimplementedChannelStreamServiceServer
	qe query.Execute
}

func NewServerRPC(qe query.Execute) *ServerRPC {
	return &ServerRPC{qe: qe}
}

func (s *ServerRPC) BindTo(srv *grpc.Server) {
	api.RegisterChannelStreamServiceServer(srv, s)
}

func (s *ServerRPC) Create(stream api.ChannelStreamService_CreateServer) error {
	ch := make(chan *models.ChannelSample)
	rfl := model.NewReflect(&ch)
	goe := tsquery.NewCreate().Model(rfl).BindExec(s.qe).GoExec(stream.Context())
	wg := errgroup.Group{}
	wg.Go(func() error { return <-goe.Errors })
	wg.Go(func() error {
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			sRfl := model.NewReflect(&models.ChannelSample{})
			rpc.NewModelExchange(req.CCR, sRfl).ToDest()
			rfl.ChanSend(sRfl)
		}
	})
	return wg.Wait()
}
