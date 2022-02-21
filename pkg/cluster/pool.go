package cluster

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"google.golang.org/grpc"
)

type NodeRPCPool struct {
	*rpc.Pool
}

func (np *NodeRPCPool) Retrieve(node *models.Node) (*grpc.ClientConn, error) {
	addr, err := node.RPCAddress()
	if err != nil {
		return nil, err
	}
	return np.Pool.Retrieve(addr)
}
