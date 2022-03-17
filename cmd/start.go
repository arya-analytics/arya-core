// Copyright © 2022 Arya Analytics

package cmd

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/api/rpc/bulktelem"
	"github.com/arya-analytics/aryacore/pkg/api/rpc/chanconfig"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/arya-analytics/aryacore/pkg/storage/minio"
	"github.com/arya-analytics/aryacore/pkg/storage/redis"
	"github.com/arya-analytics/aryacore/pkg/storage/roach"
	telemchanchunk "github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the arya core server",
	Long:  "Start the arya core server",
	Args:  cobra.NoArgs,
	RunE:  runStart,
}

func init() {
	rootCmd.AddCommand(startCmd)
	configureStartFlags()
}

func configureStartFlags() {
	flags := []configureFlag{
		configureListenAddressFlag,
	}
	for _, cf := range flags {
		cf()
	}
}

var (
	listenAddrFlag             = "listen-addr"
	listenAddressFlagShortHand = "l"
)

func configureListenAddressFlag() {
	def := fmt.Sprintf("localhost:%v", models.NodeDefaultRPCPort)
	startCmd.PersistentFlags().StringP(
		listenAddrFlag,
		listenAddressFlagShortHand,
		def,
		"address for the arya core server to listen on",
	)
	viper.SetDefault(listenAddrFlag, def)
	viper.BindPFlag(listenAddrFlag, startCmd.Flags().Lookup(listenAddrFlag))
}

func runStart(cmd *cobra.Command, _ []string) error {
	store, sErr := startStorage(cmd)

	if sErr != nil {
		return sErr
	}

	clust := startCluster(cmd, store)

	rngSvc, rngErr := startRngSvc(cmd, clust)

	if rngErr != nil {
		return rngErr
	}

	chanSvc := startChanChunkSvc(cmd, rngSvc, clust)

	if gErr := startGRPCServer(clust, chanSvc); gErr != nil {
		return gErr
	}

	for {
		select {
		case err := <-store.Errors():
			return err
		case err := <-rngSvc.Errors():
			return err
		}
	}
}

func startStorage(cmd *cobra.Command) (storage.Storage, error) {
	pool := internal.NewPool()

	mdDriver := roach.DriverRoach{Config: roach.Config{}.Viper()}
	mdEngine := roach.New(mdDriver, pool)

	objDriver := minio.DriverMinio{Config: minio.Config{}.Viper()}
	objEngine := minio.New(objDriver, pool)

	cacheDriver := redis.DriverRedis{Config: redis.Config{}.Viper()}
	cacheEngine := redis.New(cacheDriver, pool)

	s := storage.New(storage.Config{EngineMD: mdEngine, EngineObject: objEngine, EngineCache: cacheEngine})
	models.BindHooks(s)

	if err := s.NewMigrate().Exec(cmd.Context()); err != nil {
		return s, err
	}

	return s, s.Start(cmd.Context())
}

func startCluster(cmd *cobra.Command, store storage.Storage) cluster.Cluster {
	pool := startNodeRPCPool()
	clust := cluster.New()
	clust.BindService(chanchunk.NewService(store.Exec, chanchunk.NewServiceRemoteRPC(pool)))
	clust.BindService(cluster.NewStorageService(store))
	return clust
}

func startNodeRPCPool() *cluster.NodeRPCPool {
	rpcPool := rpc.NewPool(grpc.WithTransportCredentials(insecure.NewCredentials()))
	return &cluster.NodeRPCPool{Pool: rpcPool}
}

func startRngSvc(cmd *cobra.Command, clust cluster.Cluster) (*rng.Service, error) {
	obs := rng.NewObserveMem([]rng.ObservedRange{})
	if err := rng.RetrieveAddOpenRanges(cmd.Context(), clust.Exec, obs); err != nil {
		return nil, err
	}
	rngSvc := rng.NewService(obs, clust.Exec)
	rngSvc.Start(cmd.Context())
	return rngSvc, nil
}

func startChanChunkSvc(cmd *cobra.Command, rngSvc *rng.Service, clust cluster.Cluster) *telemchanchunk.Service {
	obs := telemchanchunk.NewObserveMem()
	svc := telemchanchunk.NewService(clust.Exec, obs, rngSvc)
	return svc
}

const lisNetwork = "tcp"

func startGRPCServer(clust cluster.Cluster, chanChunkSvc *telemchanchunk.Service) error {
	la := viper.GetString(listenAddrFlag)
	lis, lisErr := net.Listen(lisNetwork, la)
	log.Infof("GRPC Server listening on %s", la)
	if lisErr != nil {
		return lisErr
	}
	grpcServer := grpc.NewServer()

	// || CLUSTER CHANCHUNK ||

	persist := &chanchunk.ServerRPCPersistCluster{Cluster: clust}
	ccServer := chanchunk.NewServerRPC(persist)
	ccServer.BindTo(grpcServer)

	// || TELEM CHANCHUNK ||

	bulkTelemServer := bulktelem.NewServer(chanChunkSvc)
	bulkTelemServer.BindTo(grpcServer)

	// || TELEM CHANCONFIG ||

	chanConfigServer := chanconfig.NewServer(clust)
	chanConfigServer.BindTo(grpcServer)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalln(err)
		}
	}()
	return nil
}
