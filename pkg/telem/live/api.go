package live

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/ds"
	"github.com/arya-analytics/aryacore/pkg/server"
	"github.com/arya-analytics/aryacore/pkg/telem"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/uptrace/bun"
)

var eBases = []string{"telem", "live"}

var upgrader = websocket.Upgrader{}

func API(sc *server.Context) server.APISlice {
	eb := sc.EndpointBuilder.Child(eBases...)
	loc := NewLocator(sc.Pooler)
	rel := NewRelay(loc)
	rcv := NewReceiver(rel)
	go rcv.Start(ReceiverConfig{})
	return server.APISlice{
		OnStart: func() func() {
			go rel.Start()
			return func() {}
		},
		Handlers: []server.APIHandler{
			{
				Endpoint:    eb.Build("pull"),
				ReqMethods:  []string{"GET"},
				AuthMethods: []ds.AuthMethod{ds.Token, ds.TLS},
				HandlerFunc: func(gc *gin.Context) {
					conn, err := upgrader.Upgrade(gc.Writer, gc.Request, nil)
					if err != nil {
						panic(err)
					}
					s := NewSender(rel, conn)
					go s.start()
				},
			},
			{
				Endpoint:    eb.Build("scratch"),
				ReqMethods:  []string{"GET"},
				AuthMethods: []ds.AuthMethod{ds.Token, ds.TLS},
				HandlerFunc: func(c *gin.Context) {
					db := sc.Pooler.GetOrCreate("aryadb").(*bun.DB)
					ctx := context.Background()
					var cfgs []telem.ChannelConfig
					if err := db.NewSelect().Model(&cfgs).Scan(ctx); err != nil {
						panic(err)
					}
					fmt.Println(cfgs)
				},
			},
		},
	}

}
