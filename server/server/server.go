package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"zyrouge.me/umi/config"
	"zyrouge.me/umi/utils"
)

func Start() error {
	config, err := config.GetConfig()
	if err != nil {
		return err
	}
	server := http.NewServeMux()
	router := mux.NewRouter()
	// if err := lending.AttachRoutes(router.PathPrefix("/api/lending").Subrouter()); err != nil {
	// 	return err
	// }
	router.HandleFunc("/api/ping", PingRoute)
	AttachWebFilesRoutes(router)
	handler := http.Handler(router)
	if len(config.Server.CrossOrigin) > 0 {
		handler = CorsMiddleware(config.Server.CrossOrigin, handler)
	}
	server.Handle("/", handler)
	listenAddr := config.Server.Host + ":" + strconv.Itoa(config.Server.Port)
	utils.Logger.Info().Msg("listening on http://" + listenAddr)
	return http.ListenAndServe(listenAddr, server)
}
