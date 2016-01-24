package api

import (
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"store"
)

type Server struct {
	Store  *store.StoreManager
	Router *mux.Router
}

var Srv *Server

func NewServer() {

	glog.Info("Server is initializing...")

	Srv = &Server{}
	Srv.Store = store.NewStore()
	Srv.Router = mux.NewRouter()
	// Srv.Router.NotFoundHandler = http.HandlerFunc(Handle404)
}
