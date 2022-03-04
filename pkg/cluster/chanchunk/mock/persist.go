package mock

import (
	"context"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanchunk/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type ServerRPCPersist struct {
	ChunkReplicas []*models.ChannelChunkReplica
}

func (sp *ServerRPCPersist) RetrieveReplica(ctx context.Context, ccr *api.ChannelChunkReplica, pk model.PK) error {
	for _, mCCR := range sp.ChunkReplicas {
		if model.NewPK(mCCR.ID).Equals(pk) {
			rpc.NewModelExchange(mCCR, ccr).ToDest()
			return nil
		}
	}
	return query.Error{Type: query.ErrorTypeItemNotFound}
}
func (sp *ServerRPCPersist) CreateReplica(ctx context.Context, ccr *api.ChannelChunkReplica) error {
	mCCR := &models.ChannelChunkReplica{}
	rpc.NewModelExchange(ccr, mCCR).ToDest()
	sp.ChunkReplicas = append(sp.ChunkReplicas, mCCR)
	return nil
}

func (sp *ServerRPCPersist) DeleteReplicas(ctx context.Context, pkc model.PKChain) error {
	for i, mCCR := range sp.ChunkReplicas {
		for _, pk := range pkc {
			if model.NewPK(mCCR.ID).Equals(pk) {
				sp.ChunkReplicas = append(sp.ChunkReplicas[:i], sp.ChunkReplicas[i+1:]...)
			}
		}
	}
	return nil
}
