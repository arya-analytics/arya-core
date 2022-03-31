package chanstream

import (
	chanstreamv1 "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"io"
)

type Server struct {
	chanstreamv1.UnimplementedChannelStreamServiceServer
	svc *chanstream.Service
}

func NewServer(svc *chanstream.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) BindTo(srv *grpc.Server) {
	chanstreamv1.RegisterChannelStreamServiceServer(srv, s)
}

func (s *Server) RetrieveStream(rpcStream chanstreamv1.ChannelStreamService_RetrieveServer) error {
	ch := make(chan *models.ChannelSample)
	var (
		wg      errgroup.Group
		qStream = &streamq.Stream{Errors: make(chan error, 1)}
	)
	wg.Go(func() error { return <-qStream.Errors })
	wg.Go(func() error {
		for {
			req, err := rpcStream.Recv()
			if err == io.EOF {
				return nil
			}
			select {
			case <-rpcStream.Context().Done():
				return nil
			default:
			}

			pkc, err := model.NewPK(uuid.UUID{}).NewChainFromStrings(req.PKC...)
			if err != nil {
				return err
			}
			ch = make(chan *models.ChannelSample, len(pkc))
			if qStream, err = streamq.
				NewTSRetrieve().
				Model(&ch).
				WherePKs(pkc).BindExec(s.svc.Exec).Stream(rpcStream.Context()); err != nil {
				return err
			}
			for s := range ch {
				res := &chanstreamv1.RetrieveResponse{CCR: &chanstreamv1.ChannelSample{}}
				exc := rpc.NewModelExchange(res.CCR, s)
				exc.ToSource()
				if err := rpcStream.Send(res); err != nil {
					return err
				}
			}
		}
	})
	return wg.Wait()
}

func (s *Server) CreateStream(rpcStream chanstreamv1.ChannelStreamService_CreateServer) error {
	ch := make(chan *models.ChannelSample)
	rfl := model.NewReflect(&ch)
	stream, err := streamq.NewTSCreate().Model(rfl).BindExec(s.svc.Exec).Stream(rpcStream.Context())
	if err != nil {
		return err
	}
	wg := errgroup.Group{}
	wg.Go(func() error { return <-stream.Errors })
	wg.Go(func() error {
		for {
			req, err := rpcStream.Recv()
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
