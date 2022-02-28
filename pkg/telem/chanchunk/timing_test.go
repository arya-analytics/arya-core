package chanchunk_test

//var _ = Describe("Timing", func() {
//	var (
//		chunkSize    int64 = 10000
//		config       *models.ChannelConfig
//		lastCC       *models.ChannelChunk
//		lastCCR      *models.ChannelChunkReplica
//		lastSampleTS telem.TimeStamp
//		obs          chanchunk.Observe
//		timing       chanchunk.Timing
//	)
//	BeforeEach(func() {
//		config = &models.ChannelConfig{
//			ID:             uuid.New(),
//			NodeID:         1,
//			DataType:       telem.DataTypeFloat64,
//			DataRate:       telem.DataRate(1),
//			ConflictPolicy: models.ChannelConflictPolicyDiscard,
//		}
//		lastCC = &models.ChannelChunk{
//			ID:              uuid.New(),
//			ChannelConfig:   config,
//			ChannelConfigID: config.ID,
//			StartTS:         telem.NewTimeStamp(time.Now()),
//			Size:            chunkSize,
//		}
//		data := telem.NewChunkData([]byte{})
//		mock.TelemBulkPopulateRandomFloat64(data, int(chunkSize))
//		lastCCR = &models.ChannelChunkReplica{
//			ID:             uuid.New(),
//			ChannelChunkID: lastCC.ID,
//			Telem:           data,
//		}
//		telemChunk := telem.NewChunk(lastCC.StartTS, config.DataType, config.DataRate, lastCCR.Telem)
//		lastSampleTS = telemChunk.End()
//		obs = chanchunk.NewObserveMem([]chanchunk.ObservedChannel{{
//			ConfigPK:       config.ID,
//			DataType:       config.DataType,
//			DataRate:       config.DataRate,
//			ConflictPolicy: config.ConflictPolicy,
//			LatestSampleTS: lastSampleTS,
//		}})
//		timing = chanchunk.Timing{Obs: obs}
//	})
//	Describe("Conflict Resolution", func() {
//		It("Should resolve a conflict between two overlapping chunks", func() {
//			td := telem.NewChunkData([]byte{})
//			mock.TelemBulkPopulateRandomFloat64(td, int(chunkSize))
//			err := timing.Exec(config.ID, lastSampleTS.Add(telem.NewTimeSpan(-1*time.Second)), td)
//			Expect(err).To(BeNil())
//			Expect(td.Size()).To(Equal(79992))
//		})
//
//	})
//})
