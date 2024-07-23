package api

import (
	"database/sql"
	"log"
	"net/http"
	"test/config"
	"test/service/chat"
	"test/service/user"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

type APIServer struct {
	addr        string
	db          *sql.DB
	session     *sessions.CookieStore
	userHandler *user.Handler
	chatHandler *chat.Handler
	upgrader    websocket.Upgrader
	wsManager   *WebSocketManager
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	session := sessions.NewCookieStore([]byte(config.Envs.SesSecret))

	userStore := user.NewStore(db)
	userHandler := user.NewHandler(userStore, session)

	chatStore := chat.NewStore(db)
	chatHandler := chat.NewHandler(chatStore, session)

	wsManager := NewWebSocketManager()

	return &APIServer{
		addr:        addr,
		db:          db,
		session:     session,
		userHandler: userHandler,
		chatHandler: chatHandler,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		wsManager: wsManager,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	s.userHandler.UserRoutes(subrouter)

	s.chatHandler.ChatRoutes(subrouter)

	subrouter.HandleFunc("/ws", s.handleWebSocket)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	go s.wsManager.Run()

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, handler)
}
