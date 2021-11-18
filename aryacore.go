package aryacore

import (
	"context"
	"github.com/arya-analytics/aryacore/config"
	"github.com/arya-analytics/aryacore/server"
	"github.com/arya-analytics/aryacore/telem/live"
)



func StartServer() {
	sv := server.New(config.GetConfig(), context.Background())
	sv.BindSlice(live.API)
	sv.Start()
}
