package chat

import (
	"fmt"
	"net/http"
	"test/types"
	"test/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Handler struct {
	store     types.ChatStore
	session   *sessions.CookieStore
	wsManager *WebSocketManager
}

func NewHandler(store types.ChatStore, session *sessions.CookieStore) *Handler {
	return &Handler{
		store:     store,
		session:   session,
		wsManager: NewWebSocketManager(),
	}
}

func (h *Handler) ChatRoutes(router *mux.Router) {
	router.HandleFunc("/to-chat", h.handleToChat).Methods("Post")
	router.HandleFunc("/chats", h.handleGetChats).Methods("GET")
	router.HandleFunc("/message", h.handleGetMessages).Methods("GET")
	router.HandleFunc("/message", h.handleSendMessage).Methods("POST")
	router.HandleFunc("/chat-users", h.handleGetChatUsers).Methods("GET")
	router.HandleFunc("/ws/{chatId}", h.wsManager.HandleWebSocket)
}

func (h *Handler) handleGetChats(w http.ResponseWriter, r *http.Request) {
	cookie, _ := h.session.Get(r, "kukis")

	userId := cookie.Values["user"]

	if userId == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not authorized"))
		return
	}

	chats := h.store.GetChatsByUserID(userId.(string))

	utils.WriteJSON(w, http.StatusOK, chats)
}

func (h *Handler) handleToChat(w http.ResponseWriter, r *http.Request) {
	var payload types.GetOrCreateChatPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	chat, _ := h.store.GetOrCreateChat(payload)

	utils.WriteJSON(w, http.StatusOK, chat)
}

func (h *Handler) handleGetMessages(w http.ResponseWriter, r *http.Request) {
	chat_id := r.URL.Query().Get("chat_id")

	messages, _ := h.store.GetMessages(chat_id)

	utils.WriteJSON(w, http.StatusOK, messages)
}

func (h *Handler) handleGetChatUsers(w http.ResponseWriter, r *http.Request) {
	chat_id := r.URL.Query().Get("chat_id")

	cookie, _ := h.session.Get(r, "kukis")

	userId := cookie.Values["user"]

	if userId == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not authorized"))
		return
	}

	users := h.store.GetChatUsers(userId.(string), chat_id)

	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *Handler) handleSendMessage(w http.ResponseWriter, r *http.Request) {
	var payload types.SendMessagePayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	cookie, _ := h.session.Get(r, "kukis")

	userId := cookie.Values["user"]

	if userId == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not authorized"))
		return
	}

	message := types.Message{
		Message:  payload.Message,
		UserID:   userId.(string),
		UserName: payload.Name,
		ChatID:   payload.ChatID,
	}

	err := h.store.CreateMessage(message)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	message.UserID = ""

	h.wsManager.BroadcastToChat(payload.ChatID, message)

	utils.WriteJSON(w, http.StatusCreated, nil)
}
