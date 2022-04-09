package chanchunk_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

var _ = Describe("StreamCreate", func() {
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
			ID:             uuid.New(),
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
			streamQ     *streamq.Stream
			chunkStream chan chanchunk.StreamCreateArgs
			cancel      context.CancelFunc
		)
		JustBeforeEach(func() {
			var aCtx context.Context
			aCtx, cancel = context.WithCancel(ctx)
			chunkStream = make(chan chanchunk.StreamCreateArgs)
			var err error
			streamQ, err = svc.NewTSCreate().WhereConfigPK(config.ID).Model(&chunkStream).Stream(aCtx)
			Expect(err).To(BeNil())
			go func() {
				defer GinkgoRecover()
				Expect(<-streamQ.Errors).To(BeNil())
			}()
		})
		Describe("The  basics", func() {
			It("Should create a single new telemetry chunk correctly", func() {
				data := telem.NewChunkData([]byte{})
				Expect(data.WriteData([]float64{1, 2, 3, 4})).To(BeNil())

				By("Sending a new chunk")
				chunkStream <- chanchunk.StreamCreateArgs{Start: telem.TimeStamp(0), Data: data}

				By("Closing the stream")
				cancel()
				streamQ.Wait()

				time.Sleep(5 * time.Millisecond)

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
				cc := mock.ChunkSet(
					5,
					telem.TimeStamp(0),
					telem.DataTypeFloat64,
					telem.DataRate(25),
					telem.NewTimeSpan(1*time.Minute),
					telem.TimeSpan(0),
				)
				for _, c := range cc {
					chunkStream <- chanchunk.StreamCreateArgs{Start: c.Start(), Data: c.ChunkData}
				}

				By("Closing the stream")
				cancel()
				streamQ.Wait()

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
				cc := mock.ChunkSet(
					5,
					telem.TimeStamp(0),
					telem.DataTypeFloat64,
					telem.DataRate(25),
					telem.NewTimeSpan(1*time.Minute),
					telem.NewTimeSpan(-1*time.Second),
				)
				for _, c := range cc {
					chunkStream <- chanchunk.StreamCreateArgs{Start: c.Start(), Data: c.ChunkData}
				}

				By("Closing the stream")
				cancel()
				streamQ.Wait()

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
					Expect(tc.Span()).To(Equal(telem.NewTimeSpan(59 * time.Second)))
					Expect(tc.Start()).To(Equal(resTC[i-1].End()))
				}
			})
			It("Should update the channel cfg states", func() {
				data := telem.NewChunkData([]byte{})
				Expect(data.WriteData([]float64{1, 2, 3, 4})).To(BeNil())

				time.Sleep(50 * time.Millisecond)
				resCCActive := &models.ChannelConfig{}
				Expect(clust.NewRetrieve().Model(resCCActive).WherePK(config.ID).Exec(ctx)).To(BeNil())
				Expect(resCCActive.ID).To(Equal(config.ID))
				Expect(resCCActive.State).To(Equal(models.ChannelStatusActive))

				By("Closing the streamq")
				cancel()
				streamQ.Wait()

				resCCInactive := &models.ChannelConfig{}
				Expect(clust.NewRetrieve().Model(resCCInactive).WherePK(config.ID).Exec(ctx)).To(BeNil())
				Expect(resCCInactive.ID).To(Equal(config.ID))
				Expect(resCCInactive.State).To(Equal(models.ChannelStatusInactive))
			})
			Describe("Opening a streamq to a channel that already has Data", func() {
				It("Should create all the chunks correctly", func() {
					cc := mock.ChunkSet(
						5,
						telem.TimeStamp(0),
						telem.DataTypeFloat64,
						telem.DataRate(25),
						telem.NewTimeSpan(1*time.Minute),
						telem.TimeSpan(0),
					)
					for i, c := range cc {
						if i == 0 {
							chunkStream <- chanchunk.StreamCreateArgs{Start: c.Start(), Data: c.ChunkData}
						}
					}

					cancel()
					streamQ.Wait()

					aCtxTwo, cancelTwo := context.WithCancel(context.Background())
					var err error
					streamQ, err = svc.NewTSCreate().WhereConfigPK(config.ID).Model(&chunkStream).Stream(aCtxTwo)
					Expect(err).To(BeNil())

					for i, c := range cc {
						if i != 0 {
							chunkStream <- chanchunk.StreamCreateArgs{Start: c.Start(), Data: c.ChunkData}
						}
					}

					cancelTwo()
					streamQ.Wait()

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
			})
		})
	})
	Describe("Edge cases + errors", func() {
		Describe("Non contiguous chunks", func() {
			var (
				streamQ     *streamq.Stream
				chunkStream chan chanchunk.StreamCreateArgs
				cancel      context.CancelFunc
			)
			JustBeforeEach(func() {
				var aCtx context.Context
				aCtx, cancel = context.WithCancel(ctx)
				chunkStream = make(chan chanchunk.StreamCreateArgs)
				var err error
				streamQ, err = svc.NewTSCreate().WhereConfigPK(config.ID).Model(&chunkStream).Stream(aCtx)
				Expect(err).To(BeNil())
			})
			Context("Chunks in reverse order", func() {
				It("Should return a non-contiguous chunk error and not save the chunk", func() {
					var errors []error
					go func() {
						for err := range streamQ.Errors {
							errors = append(errors, err)
						}
					}()

					cc := mock.ChunkSet(
						2,
						telem.TimeStamp(0),
						telem.DataTypeFloat64,
						telem.DataRate(25),
						telem.NewTimeSpan(1*time.Minute),
						telem.TimeSpan(0),
					)
					// Reversing the right direction
					for i := range cc {
						c := cc[len(cc)-1-i]
						chunkStream <- chanchunk.StreamCreateArgs{Start: c.Start(), Data: c.ChunkData}
					}

					cancel()
					streamQ.Wait()

					Expect(errors).To(HaveLen(1))
					Expect(errors[0].(chanchunk.Error).Type).To(Equal(chanchunk.ErrorTimingNonContiguous))

					var resCC []*models.ChannelChunk
					Expect(clust.NewRetrieve().
						Model(&resCC).
						WhereFields(query.WhereFields{"ChannelConfigID": config.ID}).
						Order(query.OrderASC, "StartTS").
						Exec(ctx)).To(BeNil())
					Expect(len(resCC)).To(Equal(1))
				})
			})
			Context("Duplicate Chunks", func() {
				It("Shouldn't discard the duplicate", func() {
					var errors []error
					wg := &sync.WaitGroup{}
					wg.Add(1)
					go func() {
						defer wg.Done()
						defer GinkgoRecover()
						err := <-streamQ.Errors
						Expect(err).To(BeNil())
					}()

					cc := mock.ChunkSet(
						2,
						telem.TimeStamp(0),
						telem.DataTypeFloat64,
						telem.DataRate(25),
						telem.NewTimeSpan(1*time.Minute),
						telem.TimeSpan(0),
					)

					// Sending the first chunk twicw
					for range cc {
						c := cc[0]
						chunkStream <- chanchunk.StreamCreateArgs{Start: c.Start(), Data: c.ChunkData}
					}

					cancel()
					streamQ.Wait()

					Expect(errors).To(HaveLen(0))

					var resCC []*models.ChannelChunkReplica
					Expect(clust.NewRetrieve().
						Model(&resCC).
						Relation("ChannelChunk", "StartTS", "Size").
						WhereFields(query.WhereFields{"ChannelChunk.ChannelConfigID": config.ID}).
						Order(query.OrderASC, "ChannelChunk.StartTS").
						Exec(ctx)).To(BeNil())
					Expect(len(resCC)).To(Equal(1))
					Expect(resCC[0].ChannelChunk.Size).To(Equal(int64(12000)))
				})
			})
			Context("Consumed chunk", func() {
				Context("Previous consumes next chunk", func() {
					It("Shouldn discard the consumed chunk", func() {
						var errors []error
						go func() {
							defer GinkgoRecover()
							Expect(<-streamQ.Errors).To(BeNil())
						}()
						cc := mock.ChunkSet(
							2,
							telem.TimeStamp(0),
							telem.DataTypeFloat64,
							telem.DataRate(25),
							telem.NewTimeSpan(1*time.Minute),
							telem.NewTimeSpan(-30*time.Second),
						)
						for i, c := range cc {
							// Creating a consumed chunk
							if i != 0 {
								c.RemoveFromEnd(c.Start().Add(telem.NewTimeSpan(25 * time.Second)))
							}
							chunkStream <- chanchunk.StreamCreateArgs{Start: c.Start(), Data: c.ChunkData}
						}

						close(chunkStream)
						cancel()

						Expect(errors).To(HaveLen(0))
						var resCC []*models.ChannelChunkReplica
						Expect(clust.NewRetrieve().
							Model(&resCC).
							Relation("ChannelChunk", "StartTS", "Size").
							WhereFields(query.WhereFields{"ChannelChunk.ChannelConfigID": config.ID}).
							Order(query.OrderASC, "ChannelChunk.StartTS").
							Exec(ctx)).To(BeNil())
						Expect(len(resCC)).To(Equal(1))
						Expect(resCC[0].ChannelChunk.Size).To(Equal(int64(12000)))
					})
				})
			})
			Context("Exec chunk consumes previous chunk", func() {
				It("Should return a non-contiguous chunk error", func() {
					wg := &sync.WaitGroup{}
					wg.Add(1)
					go func() {
						defer wg.Done()
						defer GinkgoRecover()
						err := <-streamQ.Errors
						Expect(err.(chanchunk.Error).Type).To(Equal(chanchunk.ErrorTimingNonContiguous))
					}()
					cc := mock.ChunkSet(
						2,
						telem.TimeStamp(0),
						telem.DataTypeFloat64,
						telem.DataRate(25),
						telem.NewTimeSpan(1*time.Minute),
						telem.NewTimeSpan(-30*time.Second),
					)
					for i, c := range cc {
						// Creating a consumed chunk
						if i == 0 {
							c.RemoveFromStart(c.Start().Add(telem.NewTimeSpan(35 * time.Second)))
						}
						chunkStream <- chanchunk.StreamCreateArgs{Start: c.Start(), Data: c.ChunkData}
					}
					close(chunkStream)
					cancel()
					wg.Wait()
				})
			})

		})
		Describe("Duplicate streams", func() {
			It("Should prevent duplicate write streams to the same channel", func() {
				chunkStream := make(chan chanchunk.StreamCreateArgs)
				_, err := svc.NewTSCreate().WhereConfigPK(config.ID).Model(&chunkStream).Stream(ctx)
				Expect(err).To(BeNil())
				_, err = svc.NewTSCreate().WhereConfigPK(config.ID).Model(&chunkStream).Stream(ctx)
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
