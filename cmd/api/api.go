package api

import (
	"database/sql"
	"log"
	"net/http"
	"test/service/user"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/cors"
)

type APIServer struct {
	addr    string
	db      *sql.DB
	session *sessions.CookieStore
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	session := sessions.NewCookieStore([]byte("your-secret-key"))

	return &APIServer{
		addr:    addr,
		db:      db,
		session: session,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore, s.session)
	userHandler.UserRoutes(subrouter)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, handler)
}
