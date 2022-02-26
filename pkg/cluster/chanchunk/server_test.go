package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanchunk/v1"
	mockcluster "github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	var (
		pst           *chanchunk.ServerRPCPersistCluster
		ccr           *api.ChannelChunkReplica
		clust         cluster.Cluster
		node          = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
		rng           = &models.Range{ID: uuid.New()}
		rr            = &models.RangeReplica{ID: uuid.New(), RangeID: rng.ID, NodeID: node.ID}
		cc            = &models.ChannelChunk{ID: uuid.New(), RangeID: rng.ID, ChannelConfigID: channelConfig.ID}
		items         = []interface{}{node, channelConfig, rng, rr, cc}
	)
	BeforeEach(func() {
		if clust == nil {
			var err error
			clust, err = mockcluster.New(ctx)
			Expect(err).To(BeNil())
		}
		pst = &chanchunk.ServerRPCPersistCluster{Cluster: clust}
		ccr = &api.ChannelChunkReplica{
			ID:             uuid.New().String(),
			RangeReplicaID: rr.ID.String(),
			ChannelChunkID: cc.ID.String(),
			Telem:          []byte("randomstring"),
		}
	})
	JustBeforeEach(func() {
		for _, item := range items {
			Expect(clust.NewCreate().Model(item).Exec(ctx)).To(BeNil())
		}
	})
	JustAfterEach(func() {
		for _, item := range items {
			Expect(clust.NewDelete().Model(item).WherePK(model.NewReflect(item).PK().Raw()).Exec(ctx)).To(BeNil())
		}
	})
	It("Should create the replica correctly", func() {
		Expect(pst.CreateReplica(ctx, ccr)).To(BeNil())
	})
	It("Should retrieve a replica correctly", func() {
		Expect(pst.CreateReplica(ctx, ccr)).To(BeNil())
		resCCR := &api.ChannelChunkReplica{}
		pk, pkErr := model.NewPK(uuid.UUID{}).NewFromString(ccr.ID)
		Expect(pkErr).To(BeNil())
		Expect(pst.RetrieveReplica(ctx, resCCR, pk))
		Expect(resCCR.ID).To(Equal(ccr.ID))
		Expect(resCCR.Telem).To(Equal(ccr.Telem))
	})
	It("Should delete a replica correctly", func() {
		Expect(pst.CreateReplica(ctx, ccr)).To(BeNil())
		pk, pkErr := model.NewPK(uuid.UUID{}).NewFromString(ccr.ID)
		Expect(pkErr).To(BeNil())
		Expect(pst.DeleteReplica(ctx, model.NewPKChain([]uuid.UUID{pk.Raw().(uuid.UUID)})))
		resCCR := &api.ChannelChunkReplica{}
		Expect(pst.RetrieveReplica(ctx, resCCR, pk)).ToNot(BeNil())
	})
})
