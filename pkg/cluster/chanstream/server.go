package chanstream

import (
	"context"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"io"
	"time"
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
	goe, err := streamq.NewTSCreate().Model(rfl).BindExec(s.qe).Stream(stream.Context())
	if err != nil {
		return err
	}
	wg := errgroup.Group{}
	wg.Go(func() error { return <-goe.Errors })
	wg.Go(func() error {
		for {
			req, rpcErr := stream.Recv()
			if rpcErr == io.EOF {
				return nil
			}
			if rpcErr != nil {
				return rpcErr
			}
			sample := model.NewReflect(&models.ChannelSample{})
			rpc.NewModelExchange(req.Sample, sample).ToDest()
			rfl.ChanSend(sample)
		}
	})
	return wg.Wait()
}

func (s *ServerRPC) Retrieve(rpcStream api.ChannelStreamService_RetrieveServer) error {
	sampleStream := make(chan *models.ChannelSample)
	r := &rpcRetrieve{
		rpcStream:    rpcStream,
		qe:           s.qe,
		qStream:      &streamq.Stream{Errors: make(chan error, 10)},
		sampleStream: &sampleStream,
		updateSig:    make(chan struct{}),
	}
	return r.stream()
}

type rpcRetrieve struct {
	rpcStream    api.ChannelStreamService_RetrieveServer
	cancelQ      context.CancelFunc
	qe           query.Execute
	qStream      *streamq.Stream
	sampleStream *chan *models.ChannelSample
	updateSig    chan struct{}
}

func (r *rpcRetrieve) stream() error {
	defer r.closeQuery()
	wg := errgroup.Group{}
	wg.Go(r.relayErrors)
	wg.Go(r.relaySamples)
	wg.Go(r.listenForUpdates)
	return wg.Wait()
}

func (r *rpcRetrieve) relaySamples() error {
	for {
		select {
		case s := <-*r.sampleStream:
			res := &api.RetrieveResponse{Sample: &api.ChannelSample{}}
			rpc.NewModelExchange(res.Sample, s).ToSource()
			err := r.rpcStream.Send(res)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		case <-r.rpcStream.Context().Done():
			return nil
		case <-r.updateSig:
			continue
		}
	}
}

func (r *rpcRetrieve) closeQuery() {
	if r.cancelQ != nil {
		r.cancelQ()
	}
}

func (r *rpcRetrieve) listenForUpdates() error {
	for {
		req, err := r.rpcStream.Recv()
		if err == io.EOF {
			break
		}
		select {
		case <-r.rpcStream.Context().Done():
			break
		default:
			if uErr := r.updateQuery(req.PKC); uErr != nil {
				r.qStream.Errors <- uErr
			}
		}
	}
	return nil
}

func (r *rpcRetrieve) updateQuery(pkcStr []string) error {
	time.Sleep(1 * time.Millisecond)
	pkc, err := model.NewPK(uuid.UUID{}).NewChainFromStrings(pkcStr...)
	if err != nil {
		return err
	}
	pSampleStream := make(chan *models.ChannelSample, len(pkc))
	ctx, cancel := context.WithCancel(context.Background())
	pqStream, err := streamq.
		NewTSRetrieve().
		Model(&pSampleStream).
		WherePKs(pkc).BindExec(r.qe).Stream(ctx)
	if err != nil {
		cancel()
		return err
	}
	r.closeQuery()
	r.sampleStream = &pSampleStream
	r.cancelQ = cancel
	r.qStream = pqStream
	r.updateSig <- struct{}{}
	return nil
}

func (r *rpcRetrieve) relayErrors() error {
	for err := range r.qStream.Errors {
		select {
		case <-r.rpcStream.Context().Done():
			return nil
		default:
			if rpcErr := r.rpcStream.Send(&api.RetrieveResponse{Error: &api.Error{Message: err.Error()}}); rpcErr != nil {
				return rpcErr
			}
		}

	}
	return nil
}
