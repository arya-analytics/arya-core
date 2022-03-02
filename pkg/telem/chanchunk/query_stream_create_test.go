package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"time"
)

var _ = Describe("QueryStreamCreate", func() {
	var (
		node   *models.Node
		config *models.ChannelConfig
		svc    *chanchunk.Service
		items  []interface{}
	)
	BeforeEach(func() {
		rngObs := rng.NewObserveMem([]rng.ObservedRange{})
		rngPst := rng.NewPersistCluster(clust)
		rngSvc := rng.NewService(rngObs, rngPst)
		svc = chanchunk.NewService(clust, rngSvc)
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
		Describe("The  basics", func() {
			It("Should create a single new telemetry chunk correctly", func() {
				By("Creating the stream")
				stream := svc.NewStreamCreate()

				By("Starting the stream")
				go stream.Start(ctx, config.ID)

				var streamError error
				go func() {
					defer GinkgoRecover()
					for err := range stream.Errors() {
						if err != io.EOF {
							Fail(err.Error())
						}
					}
				}()

				data := telem.NewChunkData([]byte{})
				Expect(data.WriteData([]float64{1, 2, 3, 4})).To(BeNil())

				By("Sending a new chunk")
				stream.Send(telem.TimeStamp(0), data)

				By("Closing the stream")
				stream.Close()

				By("Being error free")
				Expect(streamError).To(BeNil())

				By("Retrieving the chunk after creation")
				resCC := &models.ChannelChunk{}
				Expect(clust.NewRetrieve().
					Model(resCC).
					WhereFields(query.WhereFields{"StartTS": telem.TimeStamp(0)}).
					Exec(ctx)).To(BeNil())
				Expect(resCC.Size).To(Equal(int64(32)))
			})
		})
		Describe("Multiple Chunks", func() {
			It("Should create multiple contiguous chunks correctly", func() {
				By("Creating the stream")
				stream := svc.NewStreamCreate()

				By("Starting the stream")
				go stream.Start(ctx, config.ID)

				var streamError error
				go func() {
					defer GinkgoRecover()
					for err := range stream.Errors() {
						if err != io.EOF {
							Fail(err.Error())
						}
					}
				}()

				cc := mock.ChunkSet(
					5,
					telem.TimeStamp(0),
					telem.DataTypeFloat64,
					telem.DataRate(25),
					telem.NewTimeSpan(1*time.Minute),
					telem.TimeSpan(0),
				)
				for _, c := range cc {
					stream.Send(c.Start(), c.ChunkData)
				}

				By("Closing the stream")
				stream.Close()

				By("Being error free")
				Expect(streamError).To(BeNil())

				By("Retrieving the chunk after creation")
				var resCC []*models.ChannelChunk
				Expect(clust.NewRetrieve().
					Model(&resCC).
					WhereFields(query.WhereFields{"ChannelConfigID": config.ID}).
					Order(query.OrderASC, "StartTS").
					Exec(ctx)).To(BeNil())
				Expect(len(resCC)).To(Equal(5))
				Expect(resCC[0].Size).To(Equal(cc[0].Size()))
				Expect(resCC[4].StartTS).To(Equal(cc[4].Start()))
			})
			It("Should resolve issues with overlapping chunks", func() {
				By("Creating the stream")
				stream := svc.NewStreamCreate()

				By("Starting the stream")
				go stream.Start(ctx, config.ID)

				var streamError error
				go func() {
					defer GinkgoRecover()
					for err := range stream.Errors() {
						if err != io.EOF {
							Fail(err.Error())
						}
					}
				}()

				cc := mock.ChunkSet(
					5,
					telem.TimeStamp(0),
					telem.DataTypeFloat64,
					telem.DataRate(25),
					telem.NewTimeSpan(1*time.Minute),
					telem.NewTimeSpan(-1*time.Second),
				)
				for _, c := range cc {
					stream.Send(c.Start(), c.ChunkData)
				}

				By("Closing the stream")
				stream.Close()

				By("Being error free")
				Expect(streamError).To(BeNil())

				By("Retrieving the chunk after creation")
				var resCC []*models.ChannelChunkReplica
				Expect(clust.NewRetrieve().
					Model(&resCC).
					Relation("ChannelChunk", "StartTS", "Size").
					WhereFields(query.WhereFields{"ChannelChunk.ChannelConfigID": config.ID}).
					Order(query.OrderASC, "ChannelChunk.StartTS").
					Exec(ctx)).To(BeNil())
				Expect(len(resCC)).To(Equal(5))
				var resTC []*telem.Chunk
				for _, ccr := range resCC {
					resTC = append(resTC, telem.NewChunk(ccr.ChannelChunk.StartTS, config.DataType, config.DataRate, ccr.Telem))
				}
				for i, tc := range resTC {
					if i == 0 {
						Expect(tc.Start()).To(Equal(telem.TimeStamp(0)))
						continue
					}
					Expect(tc.Start()).To(Equal(resTC[i-1].End()))
					Expect(tc.Span()).To(Equal(telem.NewTimeSpan(59 * time.Second)))
				}
			})
		})
	})
})
