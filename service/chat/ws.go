package chat

import (
	"log"
	"net/http"
	"sync"
	"test/types"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	upgrader    websocket.Upgrader
	clients     map[string]map[*websocket.Conn]bool
	clientsLock sync.RWMutex
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // You might want to add more strict origin checking
			},
		},
		clients: make(map[string]map[*websocket.Conn]bool),
	}
}

func (wsm *WebSocketManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	chatId := mux.Vars(r)["chatId"]

	conn, err := wsm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	wsm.clientsLock.Lock()
	if _, ok := wsm.clients[chatId]; !ok {
		wsm.clients[chatId] = make(map[*websocket.Conn]bool)
	}
	wsm.clients[chatId][conn] = true
	wsm.clientsLock.Unlock()

	defer func() {
		wsm.clientsLock.Lock()
		delete(wsm.clients[chatId], conn)
		wsm.clientsLock.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (wsm *WebSocketManager) BroadcastToChat(chatId string, message types.Message) {
	wsm.clientsLock.RLock()
	defer wsm.clientsLock.RUnlock()

	for client := range wsm.clients[chatId] {
		err := client.WriteJSON(message)
		if err != nil {
			log.Printf("error: %v", err)
			client.Close()
			delete(wsm.clients[chatId], client)
		}
	}
}
