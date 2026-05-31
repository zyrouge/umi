package authentication

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"zyrouge.me/umi/application"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/utils"
)

type LoginRouteInputResponseMode string

const (
	LoginRouteInputResponseModeJson   LoginRouteInputResponseMode = "json"
	LoginRouteInputResponseModeCookie LoginRouteInputResponseMode = "cookie"
)

type LoginRouteInput struct {
	Username     string                      `json:"username"`
	Password     string                      `json:"password"`
	TeamId       *string                     `json:"team_id,omitempty"`
	ResponseMode LoginRouteInputResponseMode `json:"response_mode" validate:"required,oneof=json cookie"`
}

type LoginRouteCookieOutput struct {
	AccessTokenExpiresInSeconds  int `json:"access_token_expires_in_seconds"`
	RefreshTokenExpiresInSeconds int `json:"refresh_token_expires_in_seconds"`
}

type LoginRouteJsonOutput struct {
	AccessToken                  string `json:"access_token"`
	RefreshToken                 string `json:"refresh_token"`
	AccessTokenExpiresInSeconds  int    `json:"access_token_expires_in_seconds"`
	RefreshTokenExpiresInSeconds int    `json:"refresh_token_expires_in_seconds"`
}

func LoginRoute(w http.ResponseWriter, r *http.Request) {
	var input LoginRouteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	if input.Username == "" || input.Password == "" {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	config, err := application.GetConfig()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to get config")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	user, err := repository.GetUserByUsername(input.Username, config.Secret.UserEncryptionKeyBytes)
	if err != nil {
		utils.Logger.Error().Err(err).Str("username", input.Username).Msg("failed to query user")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	if user == nil {
		utils.WriteHttpJsonError(w, http.StatusUnauthorized, constants.UmiErrorCodeUnauthorized)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		utils.WriteHttpJsonError(w, http.StatusUnauthorized, constants.UmiErrorCodeUnauthorized)
		return
	}
	var teamId *string
	var memberRole *repository.UmiMemberRole
	if input.TeamId != nil {
		member, err := repository.GetMember(user.Id, *input.TeamId)
		if err != nil {
			utils.Logger.Error().Err(err).Str("userId", user.Id).Str("teamId", *input.TeamId).Msg("failed to query member")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
			return
		}
		if member == nil {
			utils.WriteHttpJsonError(w, http.StatusUnauthorized, constants.UmiErrorCodeUnauthorized)
			return
		}
		teamId = &member.TeamId
		memberRole = &member.Role
	}
	accessToken, err := GenerateAccessToken(user.Id, teamId, memberRole, config.Secret.JwtSecretBytes)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", user.Id).Str("teamId", utils.FormatStringPtr(teamId)).Msg("failed to generate access token")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	refreshToken, err := GenerateRefreshToken(user.Id)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", user.Id).Msg("failed to generate refresh token")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	accessTokenExpiresInSeconds := int(AccessTokenTTL.Seconds())
	refreshTokenExpiresInSeconds := int(RefreshTokenTTL.Seconds())
	if input.ResponseMode == LoginRouteInputResponseModeCookie {
		output := LoginRouteCookieOutput{
			AccessTokenExpiresInSeconds:  accessTokenExpiresInSeconds,
			RefreshTokenExpiresInSeconds: refreshTokenExpiresInSeconds,
		}
		SetSecureCookie(w, constants.HttpCookieAccessToken, accessToken, accessTokenExpiresInSeconds)
		SetSecureCookie(w, constants.HttpCookieRefreshToken, refreshToken, refreshTokenExpiresInSeconds)
		utils.WriteHttpJsonResponse(w, http.StatusOK, output)
		return
	}
	output := LoginRouteJsonOutput{
		AccessToken:                  accessToken,
		RefreshToken:                 refreshToken,
		AccessTokenExpiresInSeconds:  accessTokenExpiresInSeconds,
		RefreshTokenExpiresInSeconds: refreshTokenExpiresInSeconds,
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, output)
}
