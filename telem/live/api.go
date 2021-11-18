package live

import (
	"github.com/arya-analytics/aryacore/ds"
	"github.com/arya-analytics/aryacore/server"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var eBases = []string{"telem", "live"}

var upgrader = websocket.Upgrader{}

func API(c *server.Context) server.APISlice {
	eb := c.EndpointBuilder.Child(eBases...)
	loc := NewLocator(c.Pooler)
	rel := NewRelay(loc)
	rcv := NewReceiver(rel)
	go rcv.Start()
	return server.APISlice{
		OnStart: func () func () {
			go rel.Start()
			return func () {}
		},
		Handlers: []server.APIHandler{
			{
				Endpoint: eb.Build("pull"),
				ReqMethods: []string{"GET"},
				AuthMethods: []ds.AuthMethod{ds.Token, ds.TLS},
				HandlerFunc: func (c *gin.Context) {
					conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
					if err != nil {
						panic(err)
					}
					s := NewSender(rel, conn)
					go s.start()
				},
			},
		},
	}

}
