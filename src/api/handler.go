package api

import (
	"fmt"
	"github.com/golang/glog"
	"models"
	"net/http"
	"net/url"
	"time"
	"utils"
	"github.com/satori/go.uuid"
)

func ApiAppHandler(
	h func(*Context, http.ResponseWriter, *http.Request),
	licenses models.StringArray,
	flags models.StringArray,
) http.Handler {
	return &handler{h, licenses, flags, true}
}


type handler struct {
	handleFunc      func(*Context, http.ResponseWriter, *http.Request)
	allowedLicenses models.StringArray // допустимые лицензии
	requiredFlag    models.StringArray // обязательные флаги
	isApi           bool
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := NewContext()
	c.T.RequestId = uuid.NewV1().String()
	c.T.Ip = GetIpAddress(r)
	c.T.Path = r.URL.Path

	c.RouteName = GetRouteName(r)

	glog.V(2).Infof("\t[%v] %v, route_name='%v'", r.Method, r.URL.Path, c.RouteName)

	glog.Infof("flags = %v", h.requiredFlag)

	resp := Srv.Store.Osin.NewResponse()
	defer resp.Close()

	// ir := Srv.Store.Osin.HandleInfoRequest(resp, r)

	// if ir == nil {
	// 	c.Err = models.NewAppError()
	// 	c.Err.StatusCode = 401
	// }

	// Find client by id ir.AccessData.Client.GetId()
	// c.T.ClientId = c.Session.GetSession().UserId

	c.Session.Client.Id = uuid.FromStringOrNil("b4c8dd5b-852c-460a-9b4a-26109f9162a2")
	c.T.ClientId = c.Session.Client.Id.String()

	if utils.Cfg.ServiceSettings.Mode != utils.MODE_PROD {
		time.Sleep(time.Millisecond * utils.Cfg.ServiceSettings.TimeoutRequest)
	}

	if !h.isApi {
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Content-Security-Policy", "frame-ancestors none")
	} else {
		// All api responsed json
		w.Header().Set("Content-Type", "application/json")
	}

	w.Header().Set(models.HEADER_REQUEST_ID, c.T.RequestId)
	w.Header().Set(models.HEADER_VERSION_ID, fmt.Sprintf("%v.%v", utils.Version, utils.BuildDate))

	if c.Err == nil {
		// glog.Infof("access_token='%v', client_id='%v', scope='%#v'", ir.AccessData.AccessToken, ir.AccessData.Client.GetId(), ir.AccessData.Scope)
		h.handleFunc(c, w, r)
	}

	if c.Err != nil {
		if c.Err.StatusCode == 0 {
			c.Err.StatusCode = http.StatusBadRequest
		}

		c.Err.T = *c.T

		c.LogError(c.Err)

		if h.isApi {
			apiErrorResponse := models.NewResponseDTO()
			apiErrorResponse.TransformFrom(c.Err)

			w.WriteHeader(c.Err.StatusCode)
			w.Write(apiErrorResponse.ToJson())
		} else {
			if c.Err.StatusCode == http.StatusUnauthorized {
				http.Redirect(w, r, "/?redirect="+url.QueryEscape(r.URL.Path), http.StatusTemporaryRedirect)
				// http.Redirect(w, r, c.GetTeamURL()+"/?redirect="+url.QueryEscape(r.URL.Path), http.StatusTemporaryRedirect)
			} else {
				w.Write([]byte("RenderWebError"))
				// RenderWebError(c.Err, w, r)
			}
		}
	}
}
