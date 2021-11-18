package aryacore

import (
	"github.com/arya-analytics/aryacore/ds"
	"github.com/arya-analytics/aryacore/server"
	"github.com/arya-analytics/aryacore/telem/live"
)

type Config struct {
	DS ds.Configs
}

func StartServer() {
	cfg := GetConfig()
	pooler := ds.NewConnPooler(cfg.DS)
	eb := server.NewEndpointBuilder("api")
	ctx := server.Context{
		Pooler: pooler,
		EndpointBuilder: eb,
	}
	sv := server.New(&ctx)
	sv.BindSlice(live.API)
	sv.Start()
}
