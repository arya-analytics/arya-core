package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Cluster struct {
	storage *mock.Storage
	cluster.Cluster
}

func (c *Cluster) Stop() {
	c.storage.Stop()
}

func New(ctx context.Context, opts ...mock.StorageOpt) (*Cluster, error) {
	s := mock.NewStorage(opts...)
	if err := s.NewMigrate().Exec(ctx); err != nil {
		return nil, err
	}
	pool := &cluster.NodeRPCPool{Pool: rpc.NewPool(grpc.WithTransportCredentials(insecure.NewCredentials()))}
	baseCluster := cluster.New()
	baseCluster.BindService(chanchunk.NewService(chanchunk.NewServiceLocalStorage(s), chanchunk.NewServiceRemoteRPC(pool)))
	baseCluster.BindService(cluster.NewStorageService(s))
	c := &Cluster{storage: s, Cluster: baseCluster}
	return c, nil
}
