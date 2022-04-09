package bulktelem_test

import (
	"github.com/arya-analytics/aryacore/pkg/api/rpc/bulktelem"
	bulktelemv1 "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/bulktelem/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"net"
	"sync"
	"time"
)

var _ = Describe("Server", func() {
	var (
		node   *models.Node
		config *models.ChannelConfig
		svc    *chanchunk.Service
		cl     bulktelemv1.BulkTelemServiceClient
		items  []interface{}
	)
	BeforeEach(func() {
		// || MODEL DEFINITIONS ||
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

		// || SERVER ||
		lis, lisErr := net.Listen("tcp", "localhost:0")
		Expect(lisErr).To(BeNil())
		grpcServer := grpc.NewServer()
		server := bulktelem.NewServer(svc)
		server.BindTo(grpcServer)
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalln(err)
			}
		}()
		conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).To(BeNil())
		cl = bulktelemv1.NewBulkTelemServiceClient(conn)
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
	Describe("createStream", func() {
		It("Should create the chunks correctly", func() {
			cc := mock.ChunkSet(5, telem.TimeStamp(0), config.DataType, config.DataRate, telem.NewTimeSpan(600*time.Second), telem.TimeSpan(0))
			stream, err := cl.CreateStream(ctx)
			Expect(err).To(BeNil())

			wg := &sync.WaitGroup{}
			var errors []*bulktelemv1.Error
			wg.Add(1)
			go func() {
				for {
					res, err := stream.Recv()
					if err == io.EOF {
						wg.Done()
						break
					}
					errors = append(errors, res.Error)
				}
			}()

			for _, c := range cc {
				Expect(stream.Send(&bulktelemv1.CreateStreamRequest{
					ChannelConfigId: config.ID.String(),
					StartTs:         int64(c.Start()),
					Data:            c.Bytes(),
				})).To(BeNil())
			}

			Expect(stream.CloseSend()).To(BeNil())

			wg.Wait()

			By("Not returning any errors")
			Expect(errors).To(HaveLen(0))

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
	Describe("retrieveStream", func() {
		JustBeforeEach(func() {
			cc := mock.ChunkSet(5, telem.TimeStamp(0), config.DataType, config.DataRate, telem.NewTimeSpan(1*time.Minute), telem.TimeSpan(0))
			stream, err := cl.CreateStream(ctx)
			Expect(err).To(BeNil())

			var errors []*bulktelemv1.Error
			wg := &sync.WaitGroup{}
			wg.Add(1)
			go func() {
				for {
					res, err := stream.Recv()
					if err == io.EOF {
						wg.Done()
						break
					}
					errors = append(errors, res.Error)
				}
			}()

			for _, c := range cc {
				Expect(stream.Send(&bulktelemv1.CreateStreamRequest{
					ChannelConfigId: config.ID.String(),
					StartTs:         int64(c.Start()),
					Data:            c.Bytes(),
				})).To(BeNil())
			}
			Expect(stream.CloseSend()).To(Succeed())
			wg.Wait()
		})
		It("Should retrieve the chunks correctly", func() {
			req := &bulktelemv1.RetrieveStreamRequest{
				ChannelConfigId: model.NewPK(config.ID).String(),
				StartTs:         int64(telem.TimeStamp(0)),
				EndTs:           int64(telem.TimeStamp(0).Add(telem.NewTimeSpan(180 * time.Second))),
			}
			stream, err := cl.RetrieveStream(ctx, req)
			Expect(err).To(BeNil())
			var chunks []*bulktelemv1.RetrieveStreamResponse
			for {
				res, err := stream.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).To(BeNil())
				chunks = append(chunks, res)
			}
			Expect(chunks).To(HaveLen(4))
		})
	})
})
