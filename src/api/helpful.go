package api

import (
	"github.com/gorilla/mux"
	"models"
	"net"
	"net/http"
)

func GetIpAddress(r *http.Request) string {
	address := r.Header.Get(models.HEADER_FORWARDED)

	if len(address) == 0 {
		address = r.Header.Get(models.HEADER_REAL_IP)
	}

	if len(address) == 0 {
		address, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	return address
}

func GetRouteName(r *http.Request) string {
	if route := mux.CurrentRoute(r); route != nil {
		return route.GetName()
	}

	return ""
}
