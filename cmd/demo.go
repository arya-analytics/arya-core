/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	bulktelemv1 "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/bulktelem/v1"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// demoCmd represents the demo command
var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Runs a demonstration RPC workload on the cluster",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.Dial("localhost:26258", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalln(err)
		}
		client := bulktelemv1.NewBulkTelemServiceClient(conn)
		cc := mock.ChunkSet(200, telem.TimeStamp(0), telem.DataTypeFloat64, telem.DataRate(25), telem.NewTimeSpan(200*time.Minute), telem.TimeSpan(0))
		stream, err := client.CreateStream(cmd.Context())
		if err != nil {
			log.Fatalln(err)
		}

		wg := sync.WaitGroup{}
		wg.Add(1)
		var errors []*bulktelemv1.Error
		go func() {
			for {
				res, err := stream.Recv()
				if err == io.EOF {
					wg.Done()
					break
				}
				log.Info(res)
				log.Error(err)
				errors = append(errors, res.Error)
			}
		}()

		t0 := time.Now()
		for _, c := range cc {
			stream.Send(&bulktelemv1.CreateStreamRequest{
				ChannelConfigId: "5110fc8e-bfde-47d3-9d30-dc6c42ab6a2a",
				StartTs:         int64(c.Start()),
				Data:            c.Bytes(),
			})
		}
		log.Infof("Wrote %v samples in %v", 25*200*60*35, time.Now().Sub(t0))
		stream.CloseSend()
		wg.Wait()
		log.Info("done")
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
}
