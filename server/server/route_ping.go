package server

import (
	"net/http"

	"zyrouge.me/umi/utils"
)

func PingRoute(w http.ResponseWriter, r *http.Request) {
	utils.WriteHttpJsonResponse(w, http.StatusOK, "pong!")
}
