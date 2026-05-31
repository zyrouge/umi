package events_live

import (
	"sync"

	"zyrouge.me/umi/repository"
)

type WebsocketManager struct {
	Mutex             sync.RWMutex
	Clients           map[*WebsocketClient]struct{}
	RegisterChannel   chan *WebsocketClient
	UnregisterChannel chan *WebsocketClient
	PublishChannel    chan *repository.UmiEvent
}

var Manager = NewWebsocketManager()

func NewWebsocketManager() *WebsocketManager {
	manager := WebsocketManager{
		Clients:           make(map[*WebsocketClient]struct{}),
		PublishChannel:    make(chan *repository.UmiEvent, 256),
		RegisterChannel:   make(chan *WebsocketClient, 64),
		UnregisterChannel: make(chan *WebsocketClient, 64),
	}
	return &manager
}

func (manager *WebsocketManager) Run() {
	for {
		select {
		case client := <-manager.RegisterChannel:
			manager.Mutex.Lock()
			manager.Clients[client] = struct{}{}
			manager.Mutex.Unlock()
		case client := <-manager.UnregisterChannel:
			manager.Mutex.Lock()
			if _, ok := manager.Clients[client]; ok {
				delete(manager.Clients, client)
				close(client.Channel)
			}
			manager.Mutex.Unlock()
		case event := <-manager.PublishChannel:
			manager.Mutex.RLock()
			for client := range manager.Clients {
				if _, ok := client.EventChannels[event.ChannelId]; ok {
					select {
					case client.Channel <- event:
					default:
						go manager.Unregister(client)
					}
				}
			}
			manager.Mutex.RUnlock()
		}
	}
}

func (manager *WebsocketManager) Register(client *WebsocketClient) {
	manager.RegisterChannel <- client
}

func (manager *WebsocketManager) Unregister(client *WebsocketClient) {
	manager.UnregisterChannel <- client
}

func (manager *WebsocketManager) Publish(event *repository.UmiEvent) {
	manager.PublishChannel <- event
}

func StartWebsocketManager() error {
	go Manager.Run()
	return nil
}
