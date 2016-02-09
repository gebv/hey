package api

func InitApi() {
	r := Srv.Router.PathPrefix("/api/v1").Subrouter()

	// Demo
	InitEmpty(r)

	// Oauth
	InitOAuth2(r)

	// Components
	InitUsers(r)
	InitEvents(r)
}