package teams

import (
	"encoding/json"
	"net/http"
	"time"

	"zyrouge.me/umi/application"
	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/route_data"
	"zyrouge.me/umi/utils"
)

type CreateTeamRequest struct {
	Name string `json:"name" validate:"required,min=1,max=128"`
}

func CreateTeamRoute(w http.ResponseWriter, r *http.Request) {
	userId := authentication.GetUserId(r.Context())
	var req CreateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	now := time.Now().Unix()
	config, err := application.GetConfig()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to get config")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	teamKey, err := utils.GenerateAESKey()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to generate team key")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	encryptedKey, err := repository.EncryptTeamKey(teamKey, config.Secret.TeamEncryptionKeyBytes)
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to encrypt team key")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	teamId, err := utils.GenerateUUIDv7()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to generate team id")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	team := repository.UmiTeam{Id: teamId, Name: req.Name, EncryptionKey: encryptedKey, CreatedAt: now, UpdatedAt: now}
	if err := repository.CreateTeam(&team); err != nil {
		utils.Logger.Error().Msg("failed to create team")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	owner := repository.UmiMember{UserId: userId, TeamId: team.Id, Role: repository.UmiMemberRoleOwner, CreatedAt: now, UpdatedAt: now}
	if err := repository.InsertMember(&owner); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", team.Id).Msg("failed to create owner member")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusCreated, &team)
}

func ListTeamsRoute(w http.ResponseWriter, r *http.Request) {
	userId := authentication.GetUserId(r.Context())
	teams, err := repository.ListTeamsByUserId(userId)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", userId).Msg("failed to list teams")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, teams)
}

func GetTeamRoute(w http.ResponseWriter, r *http.Request) {
	teamId := route_data.GetTeamId(r.Context())
	team, err := repository.GetTeamById(teamId)
	if err != nil || team == nil {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.UmiErrorCodeNotFound)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, team)
}

type UpdateTeamRequest struct {
	Name string `json:"name" validate:"required,min=1,max=128"`
}

func UpdateTeamRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.UmiMemberRoleAdmin) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	var req UpdateTeamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	if err := repository.UpdateTeamName(teamId, req.Name); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Msg("failed to update team name")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	team, _ := repository.GetTeamById(teamId)
	utils.WriteHttpJsonResponse(w, http.StatusOK, team)
}

func DeleteTeamRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.UmiMemberRoleOwner) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	if err := repository.DeleteTeam(teamId); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Msg("failed to delete team")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, nil)
}
