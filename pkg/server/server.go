package server

import (
	"github.com/arya-analytics/aryacore/pkg/config"
	"github.com/arya-analytics/aryacore/pkg/ds"
	"github.com/gin-gonic/gin"
)

// || CONTEXT ||

type Context struct {
	Pooler          *ds.ConnPooler
	EndpointBuilder *EndpointBuilder
}

// || API ||

type APISlice struct {
	OnStart     func() func()
	AuthMethods []ds.AuthMethod
	Handlers    []APIHandler
}

type APIHandler struct {
	Endpoint    string
	ReqMethods  []string
	AuthMethods []ds.AuthMethod
	HandlerFunc gin.HandlerFunc
}

type APISliceFactory func(c *Context) APISlice

type Server struct {
	Context *Context
	config  *config.Config
	router  *gin.Engine
	slices  []APISlice
}

func New(cfg *config.Config) *Server {
	rtr := gin.Default()
	pooler := ds.NewConnPooler(cfg.DS)
	eb := NewEndpointBuilder(cfg.BaseEndpoint...)
	ctx := &Context{
		pooler,
		eb,
	}
	var slices []APISlice
	return &Server{ctx, cfg, rtr, slices}
}

func (sv *Server) Start() {
	var onEndFuncs []func()
	for _, slc := range sv.slices {
		onEndFuncs = append(onEndFuncs, slc.OnStart())
	}
	defer func() {
		for _, f := range onEndFuncs {
			f()
		}
	}()
	sv.router.Run()
}

func (sv *Server) BindSlice(slcFac APISliceFactory) {
	slc := slcFac(sv.Context)
	for _, h := range slc.Handlers {
		for _, m := range h.ReqMethods {
			sv.router.Handle(m, h.Endpoint, h.HandlerFunc)
		}
	}
	sv.slices = append(sv.slices, slc)
}
