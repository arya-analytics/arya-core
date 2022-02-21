package mock

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewNodeRPCPool() *cluster.NodeRPCPool {
	return &cluster.NodeRPCPool{Pool: rpc.NewPool(grpc.WithTransportCredentials(insecure.NewCredentials()))}

}
