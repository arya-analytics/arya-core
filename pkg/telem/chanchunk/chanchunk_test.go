package chanchunk_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	mockTlm "github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var ctx = context.Background()

var _ = Describe("Service", func() {
	Describe("CreateStream", func() {
		It("Should allow the creation of a stream of channel chunk replicas", func() {
			// setup
			var err error
			clust, err := mock.New(ctx)
			Expect(err).To(BeNil())
			obs := rng.NewObserveMem([]rng.ObservedRange{})
			p := &rng.PersistCluster{Cluster: clust}
			rngSVC := rng.NewService(obs, p)
			svc := chanchunk.NewService(clust, rngSVC)

			node := &models.Node{ID: 1}
			cc := &models.ChannelConfig{NodeID: node.ID}
			Expect(clust.NewCreate().Model(node).Exec(ctx)).To(BeNil())
			Expect(clust.NewCreate().Model(cc).Exec(ctx)).To(BeNil())

			stream, errChan := svc.CreateStream(ctx, cc)
			defer close(stream)
			var streamErr error
			go func() {
				streamErr = <-errChan
			}()
			var chunkReplicaIDs []uuid.UUID
			for i := 0; i < 5; i++ {
				tlm := telem.NewBulk([]byte{})
				mockTlm.TelemBulkPopulateRandomFloat64(tlm, 100)
				ccr := &models.ChannelChunkReplica{
					ID:    uuid.New(),
					Telem: tlm,
				}
				chunkReplicaIDs = append(chunkReplicaIDs, ccr.ID)
				stream <- ccr
			}
			Expect(streamErr).To(BeNil())
			var repl []*models.ChannelChunkReplica
			Expect(clust.NewRetrieve().Model(&repl).WherePKs(chunkReplicaIDs).Exec(ctx)).To(BeNil())
			Expect(repl).To(HaveLen(5))
		})
	})
})
