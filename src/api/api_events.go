package api

import (
	"github.com/gorilla/mux"
	"models"
	"net/http"
	"store"
)

func InitEvents(r *mux.Router) {
	sr := r.PathPrefix("/events").Subrouter()

	flags := models.StringArray{}

	sr.Handle("/", ApiAppHandler(CreateEvent, flags, flags)).Methods("POST")
	sr.Handle("/{event_id}", ApiAppHandler(GetEvent, flags, flags)).Methods("GET")
	// sr.Handle("/{event_id}", ApiAppHandler(UpdateUser, flags, flags)).Methods("PUT")
	// sr.Handle("/{event_id}", ApiAppHandler(EmptyHandler, flags, flags)).Methods("DELETE")
}

func CreateEvent(c *Context, w http.ResponseWriter, r *http.Request) {
	// Thread
	// Creator
	// ParentEventId
	// DataBase64 e206INC/0YDQuNCy0LXRgiDQvNC40YB9Cg== // $echo {"m": "привет мир"} | base64
	
	dto := models.NewEventDTO()
	dto.ClientId = c.Session.Client.Id.String()

	if err := dto.FromJson(r.Body); err != nil {
		c.Err = models.NewAppError()
		c.Err.Message = models.ErrNotValid.Error()
		c.Err.DevMessage = err.Error()
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	res, err := Srv.Store.Get("event").(*store.EventStore).Create(dto)

	if err != nil {
		c.Err = models.NewAppError()
		c.Err.Message = models.ErrUnknown.Error()
		c.Err.DevMessage = err.Error()
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	response := models.NewResponseDTO()
	response.StatusCode = http.StatusOK
	response.Data = res.PrimaryValue()
	w.Write(response.ToJson())
}

func GetEvent(c *Context, w http.ResponseWriter, r *http.Request) {
	dto := models.NewEventDTO()
	dto.EventId = mux.Vars(r)["event_id"]
	dto.ClientId = c.Session.Client.Id.String()

	res, err := Srv.Store.Get("event").(*store.EventStore).Create(dto)

	if err != nil {
		c.Err = models.NewAppError()
		c.Err.Message = models.ErrUnknown.Error()
		c.Err.DevMessage = err.Error()
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	response := models.NewResponseDTO()
	response.StatusCode = http.StatusOK
	response.Data = res
	w.Write(response.ToJson())
}