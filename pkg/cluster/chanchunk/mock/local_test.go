package mock_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Local", func() {
	var (
		ctx = context.Background()
		l   *mock.Local
	)
	BeforeEach(func() {
		l = mock.NewPrepopulatedLocal()
	})
	It("Should create a replica correctly", func() {
		ccr := []*models.ChannelChunkReplica{
			{
				ID:             uuid.New(),
				ChannelChunkID: l.Chunks[0].ID,
				RangeReplicaID: l.RangeReplicas[0].ID,
				Telem:          telem.NewBulk([]byte{}),
			},
		}
		ccrLen := len(l.ChunkReplicas)
		Expect(l.CreateReplica(ctx, &ccr)).To(BeNil())
		Expect(len(l.ChunkReplicas)).To(Equal(ccrLen + 1))
	})
	Describe("Retrieve Replica", func() {
		Context("By PK", func() {
			It("Should retrieve a replica correctly", func() {
				ccr := &models.ChannelChunkReplica{}
				Expect(l.RetrieveReplica(ctx, ccr, chanchunk.LocalReplicaRetrieveOpts{
					PKC: model.NewReflect(l.ChunkReplicas[0]).PKChain(),
				})).To(BeNil())
				Expect(ccr.ID).To(Equal(l.ChunkReplicas[0].ID))
			})
		})
		Context("By Fields", func() {
			It("Should retrieve a replica correctly", func() {
				var ccr []*models.ChannelChunkReplica
				Expect(l.RetrieveReplica(ctx, &ccr, chanchunk.LocalReplicaRetrieveOpts{
					WhereFields: model.WhereFields{
						"RangeReplica.NodeID": l.Nodes[0].ID,
					},
				})).To(BeNil())
				Expect(len(ccr)).To(BeNumerically(">", 0))
			})
		})
	})
	Describe("Delete Replica", func() {
		It("Should delete the replica correctly", func() {
			ccrLen := len(l.ChunkReplicas)
			Expect(l.DeleteReplica(ctx, chanchunk.LocalReplicaDeleteOpts{PKC: model.NewReflect(l.ChunkReplicas[0]).PKChain()}))
			Expect(len(l.ChunkReplicas)).To(Equal(ccrLen - 1))
		})
	})
	Describe("Update Replica", func() {
		It("Should update the replica correctly", func() {
			repl := l.ChunkReplicas[0]
			Expect(l.UpdateReplica(ctx, repl, chanchunk.LocalReplicaUpdateOpts{
				PK: model.NewReflect(repl).PK(),
			})).To(BeNil())
		})
	})
	Describe("Retrieve Range Replica", func() {
		It("Should retrieve the range replica correctly", func() {
			rr := &models.RangeReplica{}
			Expect(l.RetrieveRangeReplica(ctx, rr, chanchunk.LocalRangeReplicaRetrieveOpts{PKC: model.NewReflect(l.RangeReplicas[0]).PKChain()})).To(BeNil())
			Expect(rr.ID).To(Equal(l.RangeReplicas[0].ID))
		})
	})
})
