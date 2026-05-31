package server

import (
	"github.com/gorilla/mux"
	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/events"
	"zyrouge.me/umi/events_live"
	"zyrouge.me/umi/teams"
	"zyrouge.me/umi/users"
)

func AttachRoutes(router *mux.Router) error {
	if err := authentication.AttachRoutes(router.PathPrefix("/api/auth").Subrouter()); err != nil {
		return err
	}
	if err := users.AttachRoutes(router.PathPrefix("/api/users").Subrouter()); err != nil {
		return err
	}
	if err := teams.AttachRoutes(router.PathPrefix("/api/teams").Subrouter()); err != nil {
		return err
	}
	if err := events.AttachRoutes(router.PathPrefix("/api/events").Subrouter()); err != nil {
		return err
	}
	if err := events_live.AttachRoutes(router.PathPrefix("/api/live").Subrouter()); err != nil {
		return err
	}
	router.HandleFunc("/api/ping", PingRoute)
	return AttachWebFilesRoutes(router)
}
