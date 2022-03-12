package chanconfig_test

import (
	"github.com/arya-analytics/aryacore/pkg/api/rpc/chanconfig"
	chanconfigv1 "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/chanconfig/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

var _ = Describe("Server", func() {
	var (
		node   *models.Node
		client chanconfigv1.ChanConfigServiceClient
	)
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		lis, lisErr := net.Listen("tcp", "localhost:0")
		Expect(lisErr).To(BeNil())
		grpcServer := grpc.NewServer()
		server := chanconfig.NewServer(clust)
		server.BindTo(grpcServer)
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				log.Fatalln(err)
			}
		}()
		conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).To(BeNil())
		client = chanconfigv1.NewChanConfigServiceClient(conn)
	})
	JustBeforeEach(func() {
		Expect(clust.NewCreate().Model(node).Exec(ctx)).To(BeNil())
	})
	JustAfterEach(func() {
		Expect(clust.NewDelete().Model(node).WherePK(model.NewReflect(node).PK()).Exec(ctx)).To(BeNil())
	})
	Describe("Create Config", func() {
		It("Should create the config correctly", func() {
			By("Creating without error")
			id := uuid.New()
			config := &chanconfigv1.ChannelConfig{
				ID:             id.String(),
				NodeId:         int32(node.ID),
				Name:           "Sensor 1",
				DataType:       chanconfigv1.ChannelConfig_FLOAT64,
				DataRate:       25,
				ConflictPolicy: chanconfigv1.ChannelConfig_DISCARD,
			}
			_, err := client.CreateConfig(ctx, &chanconfigv1.CreateConfigRequest{Config: config})
			Expect(err).To(BeNil())

			By("Being able to retrieve the config after creation")
			resCC := &models.ChannelConfig{}
			Expect(clust.NewRetrieve().Model(resCC).WherePK(id).Exec(ctx)).To(BeNil())
			Expect(resCC.ID).To(Equal(id))
			Expect(resCC.DataType).To(Equal(telem.DataTypeFloat64))
			Expect(resCC.DataRate).To(Equal(telem.DataRate(25)))
			Expect(resCC.ConflictPolicy).To(Equal(models.ChannelConflictPolicyDiscard))
		})
	})
})
