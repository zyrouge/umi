package events_live

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
	"zyrouge.me/umi/application"
	"zyrouge.me/umi/authentication"
	"zyrouge.me/umi/constants"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/utils"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type LiveEventInputQuery struct {
	Channels string `schema:"channels" validate:"required"`
}

func LiveEventRoute(w http.ResponseWriter, r *http.Request) {
	var query LiveEventInputQuery
	if err := utils.GorillaSchemaDecoder.Decode(&query, r.URL.Query()); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if err := utils.GlobalValidator.Struct(query); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	var channelIds []string
	if err := json.Unmarshal([]byte(query.Channels), &channelIds); err != nil {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	if len(channelIds) == 0 {
		utils.WriteHttpJsonError(w, http.StatusBadRequest, constants.ErrorCodeInvalidInput)
		return
	}
	seen := make(map[string]struct{}, len(channelIds))
	unique := channelIds[:0]
	for _, id := range channelIds {
		if _, ok := seen[id]; !ok {
			seen[id] = struct{}{}
			unique = append(unique, id)
		}
	}
	channelIds = unique
	userId := authentication.GetUserId(r.Context())
	count, err := repository.CountAccessibleChannelsByUserIdAndChannelIds(userId, channelIds)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", userId).Strs("channelIds", channelIds).Msg("failed to check channel access")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	if count != len(channelIds) {
		utils.WriteHttpJsonError(w, http.StatusForbidden, constants.ErrorCodeForbidden)
		return
	}
	config, err := application.GetConfig()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to get master encryption key")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	channelKeys, err := repository.GetChannelTeamKeys(channelIds, config.Secret.TeamEncryptionKeyBytes)
	if err != nil {
		utils.Logger.Error().Err(err).Strs("channelIds", channelIds).Msg("failed to get channel team keys")
		utils.WriteHttpJsonError(w, http.StatusInternalServerError, constants.ErrorCodeInternal)
		return
	}
	websocketConnection, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", userId).Msg("websocket upgrade failed")
		return
	}
	id, err := utils.GenerateUUIDv7()
	if err != nil {
		utils.Logger.Error().Err(err).Msg("failed to generate client id")
		websocketConnection.Close()
		return
	}
	client := NewWebsocketClient(id, Manager, websocketConnection, channelIds)
	Manager.Register(client)
	history, err := repository.ListEventsByChannelIds(channelIds, 20, channelKeys)
	if err != nil {
		utils.Logger.Error().Err(err).Str("userId", userId).Strs("channelIds", channelIds).Msg("failed to query event history")
	} else {
		for i := len(history) - 1; i >= 0; i-- {
			client.SendEvent(history[i])
		}
	}
	go client.WritePump()
	client.ReadPump()
}
