package users

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"zyrouge.me/umi/application"
	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/utils"
)

type PatchMeRequest struct {
	Username *string `json:"username" validate:"omitempty,min=3,max=64"`
	Email    *string `json:"email" validate:"omitempty,email"`
	Password *string `json:"password" validate:"omitempty,min=8"`
}

func PatchMeRoute(w http.ResponseWriter, r *http.Request) {
	userId := authentication.GetUserId(r.Context())
	config, err := application.GetConfig()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to get config")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	userKey := config.Secret.UserEncryptionKeyBytes
	var req PatchMeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if req.Username != nil {
		existing, err := repository.GetUserByUsername(*req.Username, userKey)
		if err != nil {
			utils.Logger.Error().Err(err).Str("userId", userId).Msg("failed to query user by username")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
			return
		}
		if existing != nil && existing.Id != userId {
			utils.WriteHttpJsonError(w, http.StatusConflict, constants.ErrorCodeConflict)
			return
		}
		if err := repository.UpdateUserUsername(userId, *req.Username); err != nil {
			utils.Logger.Error().Err(err).Str("userId", userId).Msg("failed to update username")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
			return
		}
	}
	if req.Email != nil {
		if err := repository.UpdateUserEmail(userId, *req.Email, userKey); err != nil {
			utils.Logger.Error().Err(err).Str("userId", userId).Msg("failed to update email")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
			return
		}
	}
	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.Logger.Error().Err(err).Str("userId", userId).Msg("failed to hash password")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
			return
		}
		if err := repository.UpdateUserPasswordHash(userId, string(hash)); err != nil {
			utils.Logger.Error().Err(err).Str("userId", userId).Msg("failed to update password hash")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
			return
		}
	}
	user, err := repository.GetUserById(userId, userKey)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", userId).Msg("failed to get user")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, user)
}
