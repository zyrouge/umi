package users

import (
	"net/http"

	"github.com/gorilla/mux"
	"zyrouge.me/umi/authentication"
)

func AttachRoutes(router *mux.Router) error {
	r := router.NewRoute().Subrouter()
	r.Use(authentication.AuthMiddleware)
	r.Path("/me").Methods(http.MethodGet).HandlerFunc(GetMeRoute)
	r.Path("/me").Methods(http.MethodPut).HandlerFunc(PatchMeRoute)
	return nil
}
