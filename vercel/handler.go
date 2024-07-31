package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"test/db"
	"test/service/chat"
	"test/service/user"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/rs/cors"
)

var (
	server *APIServer
)

func init() {
	_db, err := db.NewPostgresStorage("postgres://default:Ut4uNix0wdRk@ep-polished-sea-a1efivnq.ap-southeast-1.aws.neon.tech:5432/verceldb?sslmode=require")
	if err != nil {
		log.Fatal(err)
	}

	db.InitStorage(_db)

	server = NewAPIServer(":4000", _db)
}

// Handler is the entry point for Vercel
func Handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

type APIServer struct {
	addr        string
	db          *sql.DB
	session     *sessions.CookieStore
	userHandler *user.Handler
	chatHandler *chat.Handler
	router      *mux.Router
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	session := sessions.NewCookieStore([]byte("515151"))

	userStore := user.NewStore(db)
	userHandler := user.NewHandler(userStore, session)

	chatStore := chat.NewStore(db)
	chatHandler := chat.NewHandler(chatStore, session)

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userHandler.UserRoutes(subrouter)
	chatHandler.ChatRoutes(subrouter)

	router.HandleFunc("/hello", userHandler.Hello).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://go-react-api-web.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	return &APIServer{
		addr:        addr,
		db:          db,
		session:     session,
		userHandler: userHandler,
		chatHandler: chatHandler,
		router:      handler.(*mux.Router),
	}
}
