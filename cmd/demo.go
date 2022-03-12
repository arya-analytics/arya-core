/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"context"
	"fmt"
	bulktelemv1 "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/bulktelem/v1"
	chanconfigv1 "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/chanconfig/v1"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"io"
	"sync"
	"time"
)

// demoCmd represents the demo command
var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Runs a demonstration RPC workload on the cluster",
}

var (
	nodeAddr       string
	configCount    int
	configNode     int
	configDataRate float64
)

func init() {
	rootCmd.AddCommand(demoCmd)
	configureDemoCmd()
}

func configureDemoCmd() {
	demoCmd.AddCommand(demoCreateCfgCmd)
	demoCmd.AddCommand(demoWriteDataCmd)
	demoCmd.PersistentFlags().IntVar(&configNode, "node", 1, "node to create configs on")
	demoCmd.PersistentFlags().StringVar(&nodeAddr, "node-addr", "localhost:26258", "node address to connect to")
	demoCmd.PersistentFlags().IntVar(&configCount, "count", 1000, "number of configs to create")
	configureDemoCreateCfgCmd()
}

func configureDemoCreateCfgCmd() {
	demoCreateCfgCmd.Flags().Float64Var(&configDataRate, "dataRate", 1000, "data rate of configs")
}

var demoCreateCfgCmd = &cobra.Command{
	Use: "config",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial(nodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		client := chanconfigv1.NewChanConfigServiceClient(conn)
		for i := 0; i < configCount; i++ {
			config := &chanconfigv1.ChannelConfig{
				ID:             uuid.New().String(),
				NodeId:         int32(configNode),
				Name:           fmt.Sprintf("Sensor %v", i),
				DataType:       chanconfigv1.ChannelConfig_FLOAT64,
				DataRate:       configDataRate,
				ConflictPolicy: chanconfigv1.ChannelConfig_DISCARD,
			}
			if _, err := client.CreateConfig(cmd.Context(), &chanconfigv1.CreateConfigRequest{Config: config}); err != nil {
				return err
			}
		}
		return nil
	},
}

var demoWriteDataCmd = &cobra.Command{
	Use: "write",
	RunE: func(cmd *cobra.Command, args []string) error {
		conn, err := grpc.Dial(nodeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		configClient := chanconfigv1.NewChanConfigServiceClient(conn)
		cfgRes, err := configClient.RetrieveConfig(cmd.Context(), &chanconfigv1.RetrieveConfigRequest{NodeId: int32(configNode), Limit: int32(configCount)})
		if err != nil {
			return err
		}
		btClient := bulktelemv1.NewBulkTelemServiceClient(conn)
		wg := sync.WaitGroup{}

		ticker := time.NewTicker(30 * time.Second)

		go func() {
			for range ticker.C {
				log.Infof("Writing samples for %v sensors to node %v.", configCount, configNode)
			}
		}()
		for _, cfg := range cfgRes.Configs {
			wg.Add(1)
			go func(cfg *chanconfigv1.ChannelConfig) {
				defer wg.Done()
				startWriteStream(cmd.Context(), btClient, cfg)
			}(cfg)
		}
		wg.Wait()
		return nil
	},
}

const timeSpan = 30 * time.Second

func startWriteStream(ctx context.Context, client bulktelemv1.BulkTelemServiceClient, config *chanconfigv1.ChannelConfig) {
	stream, err := client.CreateStream(ctx)
	defer stream.CloseSend()
	if err != nil {
		log.Fatalln(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	var errors []*bulktelemv1.Error
	go func() {
		defer wg.Done()
		for {
			res, rErr := stream.Recv()
			if rErr == io.EOF {
				wg.Done()
				break
			}
			if rErr != nil {
				log.Fatalln(rErr)
			}
			errors = append(errors, res.Error)
		}
	}()

	ticker := time.NewTicker(timeSpan)
	c := mock.ChunkSet(1,
		telem.NewTimeStamp(time.Now()),
		telem.DataType(config.DataType),
		telem.DataRate(config.DataRate),
		telem.NewTimeSpan(timeSpan),
		telem.TimeSpan(0),
	)[0]

	for range ticker.C {
		if sErr := stream.Send(&bulktelemv1.CreateStreamRequest{
			ChannelConfigId: config.ID,
			StartTs:         int64(telem.NewTimeStamp(time.Now())),
			Data:            c.Bytes(),
		}); sErr != nil {
			log.Fatalln(sErr)
		}
	}
	wg.Wait()
}
