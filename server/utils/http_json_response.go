package utils

import (
	"encoding/json"
	"net/http"

	"zyrouge.me/umi/constants"
)

func WriteHttpJsonResponse(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	payload := map[string]any{
		"success": true,
		"data":    data,
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		Logger.Error().Err(err).Msg("failed to write JSON response")
	}
}

func WriteHttpJsonError(w http.ResponseWriter, statusCode int, code constants.UmiErrorCode) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	payload := map[string]any{
		"success": false,
		"error":   code,
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		Logger.Error().Err(err).Msg("failed to write json error response")
	}
}
