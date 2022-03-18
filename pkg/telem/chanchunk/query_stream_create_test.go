package chanchunk_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sync"
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
		rngSvc := rng.NewService(rngObs, clust.Exec)
		obs := chanchunk.NewObserveMem()
		svc = chanchunk.NewService(clust.Exec, obs, rngSvc)
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
		var stream *chanchunk.QueryStreamCreate
		JustBeforeEach(func() {
			stream = svc.NewStreamCreate()
			Expect(stream.Start(ctx, config.ID)).To(BeNil())

			go func() {
				defer GinkgoRecover()
				for err := range stream.Errors() {
					Fail(err.Error())
				}
			}()
		})
		Describe("The  basics", func() {
			It("Should create a single new telemetry chunk correctly", func() {
				data := telem.NewChunkData([]byte{})
				Expect(data.WriteData([]float64{1, 2, 3, 4})).To(BeNil())

				By("Sending a new chunk")
				stream.Send(telem.TimeStamp(0), data)

				By("Closing the stream")
				stream.Close()

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
					stream.Send(c.Start(), c.ChunkData)
				}

				By("Closing the stream")
				stream.Close()

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
					stream.Send(c.Start(), c.ChunkData)
				}

				By("Closing the stream")
				stream.Close()

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
			It("Should update the channel cfg states", func() {
				data := telem.NewChunkData([]byte{})
				Expect(data.WriteData([]float64{1, 2, 3, 4})).To(BeNil())

				time.Sleep(50 * time.Millisecond)
				resCCActive := &models.ChannelConfig{}
				Expect(clust.NewRetrieve().Model(resCCActive).WherePK(config.ID).Exec(ctx)).To(BeNil())
				Expect(resCCActive.ID).To(Equal(config.ID))
				Expect(resCCActive.Status).To(Equal(models.ChannelStatusActive))

				By("Closing the stream")
				stream.Close()

				resCCInactive := &models.ChannelConfig{}
				Expect(clust.NewRetrieve().Model(resCCInactive).WherePK(config.ID).Exec(ctx)).To(BeNil())
				Expect(resCCInactive.ID).To(Equal(config.ID))
				Expect(resCCInactive.Status).To(Equal(models.ChannelStatusInactive))
			})
			Describe("Opening a stream to a channel that already has data", func() {
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
							stream.Send(c.Start(), c.ChunkData)
						}
					}

					By("Closing the stream")
					stream.Close()

					stream = svc.NewStreamCreate()

					Expect(stream.Start(ctx, config.ID)).To(BeNil())

					for i, c := range cc {
						if i != 0 {
							stream.Send(c.Start(), c.ChunkData)
						}
					}

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
			var stream *chanchunk.QueryStreamCreate
			JustBeforeEach(func() {
				stream = svc.NewStreamCreate()
				Expect(stream.Start(ctx, config.ID)).To(BeNil())
			})
			Context("Chunks in reverse order", func() {
				It("Should return a non-contiguous chunk error and not save the chunk", func() {
					var errors []error
					wg := &sync.WaitGroup{}
					go func() {
						wg.Add(1)
						defer wg.Done()
						for err := range stream.Errors() {
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
						stream.Send(c.Start(), c.ChunkData)
					}

					stream.Close()
					wg.Wait()

					Expect(errors).To(HaveLen(1))
					Expect(errors[0].(chanchunk.TimingError).Type).To(Equal(chanchunk.TimingErrorTypeNonContiguous))

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
					go func() {
						wg.Add(1)
						defer wg.Done()
						for err := range stream.Errors() {
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

					// Sending the first chunk twicw
					for range cc {
						c := cc[0]
						stream.Send(c.Start(), c.ChunkData)
					}

					stream.Close()
					wg.Wait()

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
						wg := &sync.WaitGroup{}
						go func() {
							wg.Add(1)
							defer wg.Done()
							for err := range stream.Errors() {
								errors = append(errors, err)
							}
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
							stream.Send(c.Start(), c.ChunkData)
						}
						stream.Close()
						wg.Wait()
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
			Context("Next chunk consumes previous chunk", func() {
				It("Should return a non-contiguous chunk error", func() {
					var errors []error
					wg := &sync.WaitGroup{}
					go func() {
						wg.Add(1)
						defer wg.Done()
						for err := range stream.Errors() {
							errors = append(errors, err)
						}
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
						stream.Send(c.Start(), c.ChunkData)
					}
					stream.Close()
					wg.Wait()
					Expect(errors).To(HaveLen(1))
					Expect(errors[0].(chanchunk.TimingError).Type).To(Equal(chanchunk.TimingErrorTypeNonContiguous))
				})
			})

		})
		Describe("Duplicate streams", func() {
			var stream *chanchunk.QueryStreamCreate
			JustBeforeEach(func() {
				stream = svc.NewStreamCreate()
				Expect(stream.Start(ctx, config.ID)).To(BeNil())
			})
			It("Should prevent duplicate write streams to the same channel", func() {
				streamTwo := svc.NewStreamCreate()
				err := streamTwo.Start(ctx, config.ID)
				Expect(err).ToNot(BeNil())
			})
		})
		Describe("Context Cancellation", func() {
			It("Should handle the context cancellation with grace", func() {
				stream := svc.NewStreamCreate()
				cancelCtx, cancel := context.WithCancel(ctx)

				Expect(stream.Start(cancelCtx, config.ID)).To(BeNil())

				var errors []error
				wg := sync.WaitGroup{}
				wg.Add(1)
				go func() {
					for err := range stream.Errors() {
						errors = append(errors, err)
					}
					wg.Done()
				}()

				data := telem.NewChunkData([]byte{})
				Expect(data.WriteData([]float64{1, 2, 3, 4, 5})).To(BeNil())

				// Sending the first piece of data
				start := telem.TimeStamp(0)
				stream.Send(telem.TimeStamp(0), data)

				// Cancelling the context
				time.Sleep(50 * time.Millisecond)
				cancel()

				// Sending the second piece of data
				stream.Send(start.Add(telem.NewTimeSpan(200*time.Millisecond)), data)

				// Close the stream
				stream.Close()
				wg.Wait()

				By("Resetting the channels state to inactive")
				resCCInactive := &models.ChannelConfig{}
				Expect(clust.NewRetrieve().Model(resCCInactive).WherePK(config.ID).Exec(ctx)).To(BeNil())
				Expect(resCCInactive.ID).To(Equal(config.ID))
				Expect(resCCInactive.Status).To(Equal(models.ChannelStatusInactive))

				By("Not writing the second piece of data")
				var resCC []*models.ChannelChunk
				Expect(clust.NewRetrieve().
					Model(&resCC).
					WhereFields(query.WhereFields{"ChannelConfigID": config.ID}).
					Order(query.OrderASC, "StartTS").
					Exec(ctx)).To(BeNil())
				Expect(len(resCC)).To(Equal(1))
			})
		})
	})
})
