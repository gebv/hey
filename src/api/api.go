package api

func InitApi() {
	r := Srv.Router.PathPrefix("/api/v1").Subrouter()

	// Demo
	InitEmpty(r)

	InitOAuth2(r)
}