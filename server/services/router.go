package services

import (
	"net/http"

	"github.com/gorilla/mux"
	"zyrouge.me/umi/route_data"
)

func AttachRoutes(teamRouter *mux.Router) error {
	teamRouter.Path("/services").Methods(http.MethodPost).HandlerFunc(CreateServiceRoute)
	teamRouter.Path("/services").Methods(http.MethodGet).HandlerFunc(ListServicesRoute)

	serviceRouter := teamRouter.PathPrefix("/services/{serviceId}").Subrouter()
	serviceRouter.Use(route_data.ServiceIdMiddleware)
	serviceRouter.Path("").Methods(http.MethodGet).HandlerFunc(GetServiceRoute)
	serviceRouter.Path("").Methods(http.MethodPut).HandlerFunc(UpdateServiceRoute)
	serviceRouter.Path("/rotate-token").Methods(http.MethodPost).HandlerFunc(RotateServiceTokenRoute)
	serviceRouter.Path("").Methods(http.MethodDelete).HandlerFunc(DeleteServiceRoute)

	return nil
}
