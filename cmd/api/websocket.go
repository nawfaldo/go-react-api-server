package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

type Message struct {
	Entity string `json:"entity"`
	ID     string `json:"id"`
}

type WebSocketManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (manager *WebSocketManager) Run() {
	for {
		select {
		case client := <-manager.register:
			manager.clients[client] = true
		case client := <-manager.unregister:
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				close(client.send)
			}
		case message := <-manager.broadcast:
			for client := range manager.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(manager.clients, client)
				}
			}
		}
	}
}

func (s *APIServer) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{conn: conn, send: make(chan []byte, 256)}
	s.wsManager.register <- client

	go client.writePump(s.wsManager)
	go client.readPump(s.wsManager)
}

func (c *Client) readPump(manager *WebSocketManager) {
	defer func() {
		manager.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		manager.broadcast <- message
	}
}

func (c *Client) writePump(manager *WebSocketManager) {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (s *APIServer) BroadcastChange(entity, id string) {
	message := Message{
		Entity: entity,
		ID:     id,
	}
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Println("Error marshaling message:", err)
		return
	}
	s.wsManager.broadcast <- jsonMessage
}
