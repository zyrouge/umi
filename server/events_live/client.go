package events_live

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
	"zyrouge.me/umi/repository"
	"zyrouge.me/umi/utils"
)

const (
	WebsocketWriteWait  = 10 * time.Second
	WebsocketPongWait   = 60 * time.Second
	WebsocketPingPeriod = (WebsocketPongWait * 9) / 10
)

type WebsocketClient struct {
	Id            string
	Manager       *WebsocketManager
	Connection    *websocket.Conn
	Channel       chan *repository.UmiEvent
	EventChannels map[string]struct{}
}

func NewWebsocketClient(id string, manager *WebsocketManager, connection *websocket.Conn, channels []string) *WebsocketClient {
	client := WebsocketClient{
		Id:            id,
		Manager:       manager,
		Connection:    connection,
		Channel:       make(chan *repository.UmiEvent, 128),
		EventChannels: utils.SliceToSet(channels),
	}
	return &client
}

func (client *WebsocketClient) WritePump() {
	ticker := time.NewTicker(WebsocketPingPeriod)
	defer func() {
		ticker.Stop()
		client.Connection.Close()
	}()
	for {
		select {
		case event, ok := <-client.Channel:
			client.Connection.SetWriteDeadline(time.Now().Add(WebsocketWriteWait))
			if !ok {
				client.Connection.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			data, err := json.Marshal(event)
			if err != nil {
				utils.Logger.Error().Err(err).Str("clientId", client.Id).Msg("failed to marshal event")
				continue
			}
			if err := client.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				return
			}
		case <-ticker.C:
			client.Connection.SetWriteDeadline(time.Now().Add(WebsocketWriteWait))
			if err := client.Connection.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *WebsocketClient) ReadPump() {
	defer func() {
		client.Manager.Unregister(client)
		client.Connection.Close()
	}()
	client.Connection.SetReadLimit(512)
	client.Connection.SetReadDeadline(time.Now().Add(WebsocketPongWait))
	client.Connection.SetPongHandler(func(string) error {
		client.Connection.SetReadDeadline(time.Now().Add(WebsocketPongWait))
		return nil
	})
	for {
		_, _, err := client.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				utils.Logger.Error().Err(err).Str("clientId", client.Id).Msg("websocket closed unexpectedly")
			}
			break
		}
	}
}

func (c *WebsocketClient) SendEvent(event *repository.UmiEvent) {
	select {
	case c.Channel <- event:
	default:
	}
}
