package utils

import (
	"encoding/json"
	"net/http"
)

func WriteHttpJsonResponse(w http.ResponseWriter, statusCode int, success bool, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	payload := map[string]any{
		"success": success,
		"data":    data,
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		Logger.Error().Err(err).Msg("failed to write JSON response")
	}
}
