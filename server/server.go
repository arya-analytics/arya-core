package server

import (
	"github.com/arya-analytics/aryacore/ds"
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
	ctx    *Context
	router *gin.Engine
	slices []APISlice
}

func New(c *Context) *Server {
	r := gin.Default()
	return &Server{c, r, []APISlice{}}
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
	slc := slcFac(sv.ctx)
	for _, h := range slc.Handlers {
		for _, m := range h.ReqMethods {
			sv.router.Handle(m, h.Endpoint, h.HandlerFunc)
		}
	}
	sv.slices = append(sv.slices, slc)
}
