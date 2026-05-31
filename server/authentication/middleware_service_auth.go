package authentication

import (
	"crypto/sha256"
	"net/http"

	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/route_data"
	"zyrouge.me/umi/utils"
)

func ServiceAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ExtractBearerToken(r)
		if token == "" {
			utils.WriteHttpJsonError(w, http.StatusUnauthorized, constants.UmiErrorCodeUnauthorized)
			return
		}
		tokenHashBytes := sha256.Sum256([]byte(token))
		tokenHash := utils.BytesToHex(tokenHashBytes[:])
		service, err := repository.GetServiceByTokenHash(tokenHash)
		if err != nil {
			utils.Logger.Error().Err(err).Msg("failed to query service by token hash")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
			return
		}
		if service == nil {
			utils.WriteHttpJsonError(w, http.StatusUnauthorized, constants.UmiErrorCodeUnauthorized)
			return
		}
		r = r.WithContext(route_data.WithServiceId(r.Context(), service.Id))
		r = r.WithContext(route_data.WithServiceName(r.Context(), service.Name))
		r = r.WithContext(route_data.WithTeamId(r.Context(), service.TeamId))
		next.ServeHTTP(w, r)
	})
}
