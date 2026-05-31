package authentication

import (
	"net/http"

	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/route_data"
	"zyrouge.me/umi/utils"
)

func TeamMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		teamId := route_data.GetTeamId(r.Context())
		if teamId == "" {
			utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
			return
		}
		team, err := repository.GetTeamById(teamId)
		if err != nil {
			utils.Logger.Error().Str("teamId", teamId).Err(err).Msg("failed to query team")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
			return
		}
		if team == nil {
			utils.WriteHttpJsonError(w, http.StatusNotFound, constants.ErrorCodeNotFound)
			return
		}
		userId := route_data.GetUserId(r.Context())
		member, err := repository.GetMember(userId, teamId)
		if err != nil {
			utils.Logger.Error().Str("userId", userId).Str("teamId", teamId).Err(err).Msg("failed to get member")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
			return
		}
		if member == nil {
			utils.WriteHttpJsonError(w, http.StatusForbidden, constants.ErrorCodeForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
