package users

import (
	"net/http"

	"zyrouge.me/umi/application"
	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/utils"
)

func GetMeRoute(w http.ResponseWriter, r *http.Request) {
	userId := authentication.GetUserId(r.Context())
	config, err := application.GetConfig()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to get config")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	user, err := repository.GetUserById(userId, config.Secret.UserEncryptionKeyBytes)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", userId).Msg("failed to get user")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	if user == nil {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.UmiErrorCodeNotFound)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, user)
}
