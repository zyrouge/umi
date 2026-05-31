package authentication

import (
	"net/http"

	"zyrouge.me/umi/application"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/route_data"
	"zyrouge.me/umi/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := ExtractAccessToken(r)
		if token == "" {
			utils.WriteHttpJsonError(w, http.StatusUnauthorized, constants.ErrorCodeUnauthorized)
			return
		}
		config, err := application.GetConfig()
		if err != nil {
			utils.Logger.Error().Err(err).Msg("failed to get config")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
			return
		}
		claims, err := ValidateAccessToken(token, config.Secret.JwtSecretBytes)
		if err != nil {
			utils.WriteHttpJsonError(w, http.StatusUnauthorized, constants.ErrorCodeUnauthorized)
			return
		}
		r = r.WithContext(route_data.WithUserId(r.Context(), claims.UserId))
		if claims.TeamId != nil && claims.MemberRole != nil {
			r = r.WithContext(route_data.WithTeamId(r.Context(), *claims.TeamId))
			r = r.WithContext(route_data.WithMemberRole(r.Context(), repository.UmiMemberRole(*claims.MemberRole)))
		}
		next.ServeHTTP(w, r)
	})
}
