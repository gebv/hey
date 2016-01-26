package api

import (
	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"store"
	"net/http"
	"utils"
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

func StartServer() {
	glog.Info("Starting Server...")
	glog.Infof("Server is listening on %v", utils.Cfg.ServiceSettings.ListenAddress)

	var handler http.Handler = Srv.Router

	go func() {
		// err := manners.ListenAndServe(utils.Cfg.ServiceSettings.ListenAddress, handler)
		// if err != nil {
		// 	glog.Fatalf("Error starting server, err:%v", err)
		// 	time.Sleep(time.Second)
		// 	panic("Error starting server " + err.Error())
		// }

		http.ListenAndServe(utils.Cfg.ServiceSettings.ListenAddress, handler)
	}()
}

func StopServer() {
	glog.Info("Stopping Server...")

	// manners.Close()

	glog.Info("Server stopped")
}