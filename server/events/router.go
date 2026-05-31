package events

import (
	"net/http"

	"github.com/gorilla/mux"
	"zyrouge.me/umi/authentication"
)

func AttachRoutes(router *mux.Router) error {
	serviceRouter := router.NewRoute().Subrouter()
	serviceRouter.Use(authentication.ServiceAuthMiddleware)
	serviceRouter.Path("/publish").Methods(http.MethodPost).HandlerFunc(PublishEventRoute)
	serviceRouter.Path("/publish/discord").Methods(http.MethodPost).HandlerFunc(DiscordEventRoute)
	return nil
}
