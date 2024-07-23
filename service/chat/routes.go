package chat

import (
	"net/http"
	"test/types"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Handler struct {
	store   types.ChatStore
	session *sessions.CookieStore
}

func NewHandler(store types.ChatStore, session *sessions.CookieStore) *Handler {
	return &Handler{
		store:   store,
		session: session,
	}
}

func (h *Handler) ChatRoutes(router *mux.Router) {
	router.HandleFunc("/chat", h.handleGetChat).Methods("GET")
	router.HandleFunc("/chat", h.handleSendChat).Methods("POST")
}

func (h *Handler) handleGetChat(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleSendChat(w http.ResponseWriter, r *http.Request) {

}
