package server

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"zyrouge.me/umi/application"
	"zyrouge.me/umi/utils"
)

func StartServer() error {
	config, err := application.GetConfig()
	if err != nil {
		return err
	}
	router := mux.NewRouter()
	if err := AttachRoutes(router); err != nil {
		return err
	}
	var handler http.Handler = router
	if len(config.Server.CrossOrigin) > 0 {
		handler = CorsMiddleware(handler, config.Server.CrossOrigin)
	}
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	listenAddr := config.Server.Host + ":" + strconv.Itoa(config.Server.Port)
	utils.Logger.Info().Msg("listening on http://" + listenAddr)
	return http.ListenAndServe(listenAddr, mux)
}
