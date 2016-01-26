package api

import (
	"github.com/gorilla/mux"
	"models"
	"net/http"
	"github.com/golang/glog"
	"fmt"
)

func InitEmpty(r *mux.Router) {
	sr := r.PathPrefix("/empty").Subrouter()

	flags := models.StringArray{}

	sr.Handle("/test", ApiAppHandler(EmptyHandler, flags, flags)).Methods("GET")
}

func EmptyHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	devMessage := fmt.Sprintf("[%s]: %s", r.Method, r.RequestURI)
	glog.Info(devMessage)

	res := models.NewResponseDTO()
	res.Message = fmt.Sprintf("It's ok! Empty handler.")
	res.StatusCode = http.StatusOK
	res.DevMessage = devMessage

	if r.Method == "POST" {
		data := make(map[string]interface{})

		if err := models.FromJson(&data, r.Body); err == nil {
			res.Data = data
		}

	}

	if r.URL.Query().Get("error") == "error" {
		c.Err = models.NewAppError()
		c.Err.Message = "test"
	}

	if c.Err == nil {
		w.Write(res.ToJson())
	}
}
