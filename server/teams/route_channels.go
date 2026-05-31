package teams

import (
	"encoding/json"
	"net/http"
	"time"

	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/route_data"
	"zyrouge.me/umi/utils"
)

type CreateChannelRequest struct {
	Name string `json:"name" validate:"required,min=1,max=128"`
}

func CreateChannelRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.UmiMemberRoleAdmin) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	var req CreateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	channelId, err := utils.GenerateUUIDv7()
	if err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Msg("failed to generate channel id")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	now := time.Now().Unix()
	channel := repository.UmiChannel{
		Id: channelId, Name: req.Name, TeamId: teamId, CreatedAt: now, UpdatedAt: now,
	}
	if err := repository.CreateChannel(&channel); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Msg("failed to create channel")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusCreated, &channel)
}

func ListChannelsRoute(w http.ResponseWriter, r *http.Request) {
	teamId := route_data.GetTeamId(r.Context())
	channels, err := repository.ListChannelsByTeamId(teamId)
	if err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Msg("failed to list channels")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, channels)
}

func GetChannelRoute(w http.ResponseWriter, r *http.Request) {
	teamId := route_data.GetTeamId(r.Context())
	channelId := route_data.GetChannelId(r.Context())
	channel, err := repository.GetChannelById(channelId)
	if err != nil || channel == nil || channel.TeamId != teamId {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.UmiErrorCodeNotFound)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, channel)
}

type UpdateChannelRequest struct {
	Name string `json:"name" validate:"required,min=1,max=128"`
}

func UpdateChannelRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.UmiMemberRoleAdmin) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	channelId := route_data.GetChannelId(r.Context())
	channel, err := repository.GetChannelById(channelId)
	if err != nil || channel == nil || channel.TeamId != teamId {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.UmiErrorCodeNotFound)
		return
	}
	var req UpdateChannelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(req); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	if err := repository.UpdateChannelNameById(channelId, req.Name); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Str("channelId", channelId).Msg("failed to update channel name")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	channel.Name = req.Name
	utils.WriteHttpJsonResponse(w, http.StatusOK, channel)
}

func DeleteChannelRoute(w http.ResponseWriter, r *http.Request) {
	if !authentication.RequirePermissionMiddleware(w, r, repository.UmiMemberRoleAdmin) {
		return
	}
	teamId := route_data.GetTeamId(r.Context())
	channelId := route_data.GetChannelId(r.Context())
	channel, err := repository.GetChannelById(channelId)
	if err != nil || channel == nil || channel.TeamId != teamId {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.UmiErrorCodeNotFound)
		return
	}
	if err := repository.DeleteChannelById(channelId); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", teamId).Str("channelId", channelId).Msg("failed to delete channel")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	utils.WriteHttpJsonResponse(w, http.StatusOK, nil)
}
