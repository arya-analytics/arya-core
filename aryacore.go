package aryacore

import (
	"context"
	"github.com/arya-analytics/aryacore/ds"
	"github.com/gin-gonic/gin"
)

type Config struct {
	ds ds.Configs
}

type Core struct {
	cfg         *Config
	ConnManager *ds.ConnPooler
	router      *gin.Engine
}

func NewCore(ctx context.Context, cfg *Config) (*Core, context.Context) {
	cm := ds.NewConnPooler(cfg.ds)
	router := gin.Default()
	core := &Core{
		ConnManager: cm,
		cfg:         cfg,
		router:      router,
	}
	return core, ctx
}
