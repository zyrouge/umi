package events_live

import (
	"github.com/gorilla/mux"
	"zyrouge.me/umi/authentication"
)

func AttachRoutes(router *mux.Router) error {
	router.Use(authentication.AuthMiddleware)
	router.Path("/events").HandlerFunc(LiveEventRoute)
	return nil
}
