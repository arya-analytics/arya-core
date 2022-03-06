package chanconfig

import (
	"context"
	chanconfigv1 "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/chanconfig/v1"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"google.golang.org/grpc"
)

type Server struct {
	chanconfigv1.UnimplementedChanConfigServiceServer
	clust cluster.Cluster
}

func NewServer(clust cluster.Cluster) *Server {
	return &Server{clust: clust}
}

func (s *Server) BindTo(srv *grpc.Server) {
	chanconfigv1.RegisterChanConfigServiceServer(srv, s)
}

func (s *Server) CreateConfig(ctx context.Context, req *chanconfigv1.CreateConfigRequest) (*chanconfigv1.CreateConfigResponse, error) {
	exc := rpc.NewModelExchange(&models.ChannelConfig{}, req.Config)
	exc.ToSource()
	return &chanconfigv1.CreateConfigResponse{}, s.clust.NewCreate().Model(exc.Source()).Exec(ctx)
}
