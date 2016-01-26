package api

import (
	"github.com/gorilla/mux"
	"github.com/golang/glog"
	"github.com/RangelReale/osin"
	"net/http"
	// "models"
)

func InitOAuth2(r *mux.Router) {
	sr := r.PathPrefix("/oauth2").Subrouter()

	// flags := models.StringArray{}

	sr.HandleFunc("/authorize", AuthHandler).Methods("GET", "POST")
	sr.HandleFunc("/token", TokenHandler).Methods("GET", "POST")
	sr.HandleFunc("/me", MeHandler).Methods("GET", "POST")
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	resp := Srv.Store.Osin.NewResponse()
	defer resp.Close()

	if ir := Srv.Store.Osin.HandleInfoRequest(resp, r); ir != nil {
		// fmt.Println(r.Form.Get("fields"))

		Srv.Store.Osin.FinishInfoRequest(resp, r, ir)
	}
	osin.OutputJSON(resp, w, r)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	resp := Srv.Store.Osin.NewResponse()
	defer resp.Close()

	if ar := Srv.Store.Osin.HandleAuthorizeRequest(resp, r); ar != nil {
		// ar.Client.GetId()

		glog.Infof("client_id=%v", ar.Client.GetId())

		glog.Infof("FinishAuthorizeRequest, %v", ar)

		ar.Authorized = true
		Srv.Store.Osin.FinishAuthorizeRequest(resp, r, ar)
	}

	if resp.IsError && resp.InternalError != nil {
		glog.Errorf("osin error: %s\n", resp.InternalError)
	}

	osin.OutputJSON(resp, w, r)
}

func TokenHandler(w http.ResponseWriter, r *http.Request) {
	resp := Srv.Store.Osin.NewResponse()
	defer resp.Close()

	if ar := Srv.Store.Osin.HandleAccessRequest(resp, r); ar != nil {
		ar.Authorized = true
		Srv.Store.Osin.FinishAccessRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		glog.Errorf("osin error: %s\n", resp.InternalError)
	}
	osin.OutputJSON(resp, w, r)
}