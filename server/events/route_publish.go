package events

import (
	"encoding/json"
	"net/http"
	"time"

	"zyrouge.me/umi/application"
	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/events_live"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/utils"
)

type PublishEventInput struct {
	ChannelId string                    `json:"channel_id" validate:"required"`
	Title     string                    `json:"title" validate:"required"`
	Body      *string                   `json:"body"`
	Level     *repository.UmiEventLevel `json:"level" validate:"omitempty,oneof=info warning error critical"`
	ActionURL *string                   `json:"action_url"`
	IconURL   *string                   `json:"icon_url"`
	Tags      []string                  `json:"tags"`
	Metadata  map[string]string         `json:"metadata"`
}

func PublishEventRoute(w http.ResponseWriter, r *http.Request) {
	var input PublishEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(input); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.UmiErrorCodeInvalidInput)
		return
	}
	config, err := application.GetConfig()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to get config")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	team, err := repository.GetTeamByChannelId(input.ChannelId)
	if err != nil || team == nil {
		utils.WriteHttpJsonError(w, http.StatusNotFound, constants.UmiErrorCodeNotFound)
		return
	}
	teamKey, err := repository.DecryptTeamEncryptionKey(team, config.Secret.TeamEncryptionKeyBytes)
	if err != nil {
		utils.Logger.Error().Err(err).Str("teamId", team.Id).Msg("failed to decrypt team key")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	eventId, err := utils.GenerateUUIDv7()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to generate event id")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	event := repository.UmiEvent{
		Id:        eventId,
		ChannelId: input.ChannelId,
		Title:     input.Title,
		Body:      input.Body,
		Level:     input.Level,
		ActionURL: input.ActionURL,
		IconURL:   input.IconURL,
		ServiceId: authentication.GetServiceId(r.Context()),
		Metadata:  input.Metadata,
		CreatedAt: time.Now().Unix(),
	}
	if err := repository.InsertEvent(&event, teamKey); err != nil {
		utils.Logger.Error().Err(err).Str("channelId", input.ChannelId).Msg("failed to insert event")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	now := time.Now().Unix()
	existingTags, err := repository.GetTagByNames(team.Id, input.Tags)
	if err != nil {
		utils.Logger.Error().Str("teamId", team.Id).Strs("tags", input.Tags).Err(err).Msg("failed to get existing tags")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	var newTags []*repository.UmiTag
	for _, tagName := range input.Tags {
		if _, exists := existingTags[tagName]; exists {
			continue
		}
		tagId, err := utils.GenerateUUIDv7()
		if err != nil {
			utils.Logger.Error().Err(err).Msg("failed to generate tag id")
			utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
			return
		}
		existingTags[tagName] = &repository.UmiTag{
			Id: tagId, TeamId: team.Id, Name: tagName, CreatedAt: now, UpdatedAt: now,
		}
		newTags = append(newTags, existingTags[tagName])
	}
	if err := repository.BulkInsertTags(newTags); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", team.Id).Msg("failed to bulk insert tags")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	associations := make([]*repository.UmiEventTag, 0, len(input.Tags))
	for _, tagName := range input.Tags {
		associations = append(associations, &repository.UmiEventTag{
			EventId: eventId, TagId: existingTags[tagName].Id, CreatedAt: now,
		})
	}
	if err := repository.BulkInsertEventTags(associations); err != nil {
		utils.Logger.Error().Err(err).Str("teamId", team.Id).Strs("tags", input.Tags).Msg("failed to bulk insert event tag associations")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.UmiErrorCodeInternal)
		return
	}
	events_live.Manager.Publish(&event)
	utils.WriteHttpJsonResponse(w, http.StatusCreated, &event)
}
