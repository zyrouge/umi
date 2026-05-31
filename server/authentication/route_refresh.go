package authentication

import (
	"encoding/json"
	"net/http"

	"zyrouge.me/umi/application"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/utils"
)

type RefreshRouteInput struct {
	RefreshToken string                      `json:"refresh_token"`
	TeamId       *string                     `json:"team_id,omitempty"`
	ResponseMode LoginRouteInputResponseMode `json:"response_mode" validate:"required,oneof=json cookie"`
}

type RefreshRouteCookieOutput struct {
	AccessTokenExpiresInSeconds  int `json:"access_token_expires_in_seconds"`
	RefreshTokenExpiresInSeconds int `json:"refresh_token_expires_in_seconds"`
}

type RefreshRouteJsonOutput struct {
	AccessToken                  string `json:"access_token"`
	RefreshToken                 string `json:"refresh_token"`
	AccessTokenExpiresInSeconds  int    `json:"access_token_expires_in_seconds"`
	RefreshTokenExpiresInSeconds int    `json:"refresh_token_expires_in_seconds"`
}

func RefreshRoute(w http.ResponseWriter, r *http.Request) {
	var input RefreshRouteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	refreshToken := ""
	if input.ResponseMode == LoginRouteInputResponseModeCookie {
		var err error
		refreshToken, err = GetCookieValue(r, constants.HttpCookieRefreshToken)
		if err != nil {
			utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
			return
		}
	} else {
		refreshToken = input.RefreshToken
	}
	if refreshToken == "" {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	config, err := application.GetConfig()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to get config")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	refreshTokenData, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to validate refresh token")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	user, err := repository.GetUserById(refreshTokenData.UserId, config.Secret.UserEncryptionKeyBytes)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", refreshTokenData.UserId).Msg("failed to query user")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	if user == nil {
		utils.WriteHttpJsonError(w, http.StatusUnauthorized, constants.ErrorCodeUnauthorized)
		return
	}
	var teamId *string
	var memberRole *repository.UmiMemberRole
	if input.TeamId != nil {
		member, err := repository.GetMember(user.Id, *input.TeamId)
		if err != nil {
			utils.Logger.Error().Err(err).Str("userId", user.Id).Str("teamId", utils.FormatStringPtr(input.TeamId)).Msg("failed to query member")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
			return
		}
		if member == nil {
			utils.WriteHttpJsonError(w, http.StatusUnauthorized, constants.ErrorCodeUnauthorized)
			return
		}
		teamId = &member.TeamId
		memberRole = &member.Role
	}
	newAccessToken, err := GenerateAccessToken(user.Id, teamId, memberRole, config.Secret.JwtSecretBytes)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", user.Id).Str("teamId", utils.FormatStringPtr(teamId)).Msg("failed to generate access token")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	newRefreshToken, err := GenerateRefreshToken(user.Id)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", user.Id).Msg("failed to generate refresh token")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	accessTokenExpiresInSeconds := int(AccessTokenTTL.Seconds())
	refreshTokenExpiresInSeconds := int(RefreshTokenTTL.Seconds())
	if input.ResponseMode == LoginRouteInputResponseModeCookie {
		SetSecureCookie(w, constants.HttpCookieAccessToken, newAccessToken, accessTokenExpiresInSeconds)
		SetSecureCookie(w, constants.HttpCookieRefreshToken, newRefreshToken, refreshTokenExpiresInSeconds)
		output := RefreshRouteCookieOutput{
			AccessTokenExpiresInSeconds:  accessTokenExpiresInSeconds,
			RefreshTokenExpiresInSeconds: refreshTokenExpiresInSeconds,
		}
		utils.WriteHttpJsonResponse(w, http.StatusOK, output)
		return
	}
	output := RefreshRouteJsonOutput{
		AccessToken:                  newAccessToken,
		RefreshToken:                 newRefreshToken,
		AccessTokenExpiresInSeconds:  accessTokenExpiresInSeconds,
		RefreshTokenExpiresInSeconds: refreshTokenExpiresInSeconds,
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, output)
}
