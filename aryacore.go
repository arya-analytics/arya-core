package aryacore

import (
	"context"
	"github.com/arya-analytics/aryacore/ds"
	"github.com/gin-gonic/gin"
)

type Config struct {
	ds ds.Config
}

type Core struct {
	cfg         *Config
	ConnManager *ds.ConnManager
	router      *gin.Engine
}

func NewCore(ctx context.Context, cfg *Config) (*Core, context.Context) {
	cm := ds.NewConnManager(cfg.ds)
	router := gin.Default()
	core := &Core{
		ConnManager: cm,
		cfg:         cfg,
		router:      router,
	}
	return core, ctx
}
