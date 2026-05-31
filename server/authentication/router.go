package authentication

import "github.com/gorilla/mux"

func AttachRoutes(router *mux.Router) error {
	router.Path("/login").Methods("POST").HandlerFunc(LoginRoute)
	router.Path("/refresh").Methods("POST").HandlerFunc(RefreshRoute)
	return nil
}
