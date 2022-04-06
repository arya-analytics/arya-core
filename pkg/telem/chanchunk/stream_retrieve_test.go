package chanchunk_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

var _ = Describe("streamRetrieve", func() {
	var (
		node   *models.Node
		config *models.ChannelConfig
		svc    *chanchunk.Service
		items  []interface{}
	)
	BeforeEach(func() {
		rngObs := rng.NewObserveMem([]rng.ObservedRange{})
		rngSvc := rng.NewService(rngObs, clust.Exec)
		svc = chanchunk.NewService(clust.Exec, rngSvc)
		node = &models.Node{ID: 1}
		config = &models.ChannelConfig{
			Name:           "Awesome Channel",
			NodeID:         node.ID,
			DataRate:       telem.DataRate(25),
			DataType:       telem.DataTypeFloat64,
			ConflictPolicy: models.ChannelConflictPolicyDiscard,
		}
		items = []interface{}{node, config}
	})
	JustBeforeEach(func() {
		for _, item := range items {
			Expect(clust.NewCreate().Model(item).Exec(ctx)).To(BeNil())
		}
	})
	JustAfterEach(func() {
		for _, item := range items {
			Expect(clust.NewDelete().Model(item).WherePK(model.NewReflect(item).PK()).Exec(ctx)).To(BeNil())
		}
	})
	Describe("Standard Usage", func() {
		var (
			chunkSize int64
		)
		JustBeforeEach(func() {
			aCtx, cancel := context.WithCancel(ctx)
			chunkStream := make(chan chanchunk.StreamCreateArgs)
			streamQ, err := svc.NewTSCreate().Model(&chunkStream).Stream(aCtx)
			Expect(err).To(BeNil())
			wg := &sync.WaitGroup{}
			defer func() {
				cancel()
				close(chunkStream)
				wg.Wait()
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				defer GinkgoRecover()
				Expect(<-streamQ.Errors).To(BeNil())
			}()

			cc := mock.ChunkSet(
				5,
				telem.TimeStamp(0),
				telem.DataTypeFloat64,
				telem.DataRate(25),
				telem.NewTimeSpan(1*time.Minute),
				telem.TimeSpan(0),
			)
			chunkSize = cc[0].Size()
			for _, c := range cc {
				chunkStream <- chanchunk.StreamCreateArgs{
					Start: c.Start(),
					Data:  c.ChunkData,
				}
			}
		})
		It("Should retrieve a streamq of chunks within the correct time range", func() {
			tr := telem.NewTimeRange(telem.TimeStamp(0), telem.TimeStamp(0).Add(telem.NewTimeSpan(170*time.Second)))
			chunkStream := make(chan *telem.Chunk)
			_, err := svc.NewTSRetrieve().Model(&chunkStream).WherePK(config.ID).WhereTimeRange(tr).Stream(ctx)
			Expect(err).To(BeNil())
			var chunks []*telem.Chunk
			for c := range chunkStream {
				chunks = append(chunks, c)
			}
			Expect(chunks).To(HaveLen(3))
			Expect(chunks[0].Size()).To(Equal(chunkSize))
		})
	})
})
