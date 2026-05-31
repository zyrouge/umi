package events

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"zyrouge.me/umi/application"
	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/events_live"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/utils"
)

type DiscordEventInputQuery struct {
	ChannelId string `schema:"channel_id" validate:"required"`
}

type DiscordEventInput struct {
	Content   string                   `json:"content"`
	Username  string                   `json:"username"`
	AvatarURL string                   `json:"avatar_url"`
	Embeds    []DiscordEventInputEmbed `json:"embeds"`
}

type DiscordEventInputEmbed struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
}

func DiscordEventRoute(w http.ResponseWriter, r *http.Request) {
	var query DiscordEventInputQuery
	if err := utils.GorillaSchemaDecoder.Decode(&query, r.URL.Query()); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(query); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	var input DiscordEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	title := ""
	if len(input.Embeds) > 0 && input.Embeds[0].Title != "" {
		title = input.Embeds[0].Title
	} else if input.Username != "" {
		title = input.Username
	}
	if title == "" {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	var parts []string
	if input.Content != "" {
		parts = append(parts, input.Content)
	}
	for _, embed := range input.Embeds {
		if embed.Description != "" {
			parts = append(parts, embed.Description)
		}
	}
	var body *string
	if len(parts) > 0 {
		joined := strings.Join(parts, "\n")
		body = &joined
	}
	var actionURL *string
	if len(input.Embeds) > 0 && input.Embeds[0].URL != "" {
		u := input.Embeds[0].URL
		actionURL = &u
	}
	var iconURL *string
	if input.AvatarURL != "" {
		iconURL = &input.AvatarURL
	}
	config, err := application.GetConfig()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to get master encryption key")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	team, err := repository.GetTeamByChannelId(query.ChannelId)
	if err != nil || team == nil {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.ErrorCodeNotFound)
		return
	}
	teamKey, err := repository.DecryptTeamEncryptionKey(team, config.Secret.TeamEncryptionKeyBytes)
	if err != nil {
		utils.Logger.Error().Err(err).Str("teamId", team.Id).Msg("failed to decrypt team key")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	eventId, err := utils.GenerateUUIDv7()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to generate event id")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	event := repository.UmiEvent{
		Id:        eventId,
		ChannelId: query.ChannelId,
		Title:     title,
		Body:      body,
		ActionURL: actionURL,
		IconURL:   iconURL,
		ServiceId: authentication.GetServiceId(r.Context()),
		CreatedAt: time.Now().Unix(),
	}
	if err := repository.InsertEvent(&event, teamKey); err != nil {
		utils.Logger.Error().Err(err).Str("channelId", query.ChannelId).Msg("failed to insert discord event")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	events_live.Manager.Publish(&event)
	utils.WriteHttpJsonResponse(w, http.StatusCreated, &event)
}
