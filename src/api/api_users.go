package api

import (
	"github.com/gorilla/mux"
	"models"
	"net/http"
	"store"
)

func InitUsers(r *mux.Router) {
	sr := r.PathPrefix("/users").Subrouter()

	flags := models.StringArray{}

	sr.Handle("/", ApiAppHandler(CreateUser, flags, flags)).Methods("POST")
	sr.Handle("/", ApiAppHandler(GetUser, flags, flags)).Methods("GET")
	sr.Handle("/", ApiAppHandler(UpdateUser, flags, flags)).Methods("PUT")
	sr.Handle("/", ApiAppHandler(EmptyHandler, flags, flags)).Methods("DELETE")
}

func CreateUser(c *Context, w http.ResponseWriter, r *http.Request) {
	dto := models.NewUserDTO()
	dto.ClientId = c.Session.Client.Id.String()

	if err := dto.FromJson(r.Body); err != nil {
		c.Err = models.NewAppError()
		c.Err.Message = models.ErrNotValid.Error()
		c.Err.DevMessage = err.Error()
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	res, err := Srv.Store.Get("user").(*store.UserStore).Create(dto)

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

func GetUser(c *Context, w http.ResponseWriter, r *http.Request) {
	dto := models.NewUserDTO()
	dto.ExtId = r.URL.Query().Get("ext_id")
	dto.ClientId = c.Session.Client.Id.String()

	res, err := Srv.Store.Get("user").(*store.UserStore).GetOne(dto)

	if err != nil {
		c.Err = models.NewAppError()
		c.Err.Message = models.ErrNotValid.Error()
		c.Err.DevMessage = err.Error()
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	response := models.NewResponseDTO()
	response.StatusCode = http.StatusOK
	response.Data = res
	w.Write(response.ToJson())
}

func UpdateUser(c *Context, w http.ResponseWriter, r *http.Request) {
	dto := models.NewUserDTO()
	dto.ExtId = r.URL.Query().Get("ext_id")
	dto.ClientId = c.Session.Client.Id.String()

	if err := dto.FromJson(r.Body); err != nil {
		c.Err = models.NewAppError()
		c.Err.Message = models.ErrNotValid.Error()
		c.Err.DevMessage = err.Error()
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	res, err := Srv.Store.Get("user").(*store.UserStore).Update(dto)

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

func DeleteUser(c *Context, w http.ResponseWriter, r *http.Request) {
	dto := models.NewUserDTO()
	dto.ExtId = r.URL.Query().Get("ext_id")
	dto.ClientId = c.Session.Client.Id.String()

	if err := dto.FromJson(r.Body); err != nil {
		c.Err = models.NewAppError()
		c.Err.Message = models.ErrNotValid.Error()
		c.Err.DevMessage = err.Error()
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	err := Srv.Store.Get("user").(*store.UserStore).Delete(dto)

	if err != nil {
		c.Err = models.NewAppError()
		c.Err.Message = models.ErrUnknown.Error()
		c.Err.DevMessage = err.Error()
		c.Err.StatusCode = http.StatusBadRequest
		return
	}

	response := models.NewResponseDTO()
	response.StatusCode = http.StatusOK
	w.Write(response.ToJson())
}