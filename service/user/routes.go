package user

import (
	"fmt"
	"net/http"
	"test/types"
	"test/utils"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type Handler struct {
	store   types.UserStore
	session *sessions.CookieStore
}

func NewHandler(store types.UserStore, session *sessions.CookieStore) *Handler {
	return &Handler{
		store:   store,
		session: session,
	}
}

func (h *Handler) UserRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods("POST")
	router.HandleFunc("/register", h.handleRegister).Methods("POST")
	router.HandleFunc("/auth", h.handleAuth).Methods("GET")
	router.HandleFunc("/logout", h.handleLogout).Methods("POST")
	router.HandleFunc("/users", h.handleGetUsers).Methods("GET")
	router.HandleFunc("/user", h.handleGetUser).Methods("GET")
}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload types.RegisterUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload %v", errors))
		return
	}

	userId := uuid.New().String()

	err := h.store.CreateUser(types.User{
		ID:       userId,
		Name:     payload.Name,
		Password: payload.Password,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	cookie, _ := h.session.Get(r, "kukis")
	cookie.Options = &sessions.Options{
		MaxAge:   3600 * 24, // 1 day
		SameSite: http.SameSiteStrictMode,
	}

	cookie.Values["user"] = userId
	cookie.Save(r, w)

	utils.WriteJSON(w, http.StatusCreated, nil)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload types.LoginUserPayload
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	user, err := h.store.GetUserByName(payload.Name)
	if user == nil && err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("either name or password is wrong"))
		return
	}

	if user.Password != payload.Password {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("either name or password is wrong"))
		return
	}

	cookie, _ := h.session.Get(r, "kukis")
	cookie.Options = &sessions.Options{
		MaxAge:   3600 * 24, // 1 day
		SameSite: http.SameSiteStrictMode,
	}

	cookie.Values["user"] = user.ID
	cookie.Save(r, w)

	utils.WriteJSON(w, http.StatusAccepted, nil)
}

func (h *Handler) handleAuth(w http.ResponseWriter, r *http.Request) {
	cookie, _ := h.session.Get(r, "kukis")

	userId := cookie.Values["user"]

	if userId == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("not authorized"))
		return
	}

	user, _ := h.store.GetAuthById(userId.(string))

	utils.WriteJSON(w, http.StatusAccepted, user)
}

func (h *Handler) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := h.session.Get(r, "kukis")
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	cookie.Values["user"] = nil
	cookie.Options = &sessions.Options{
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	}

	if err := cookie.Save(r, w); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil)
}

func (h *Handler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, _ := h.store.GetUsers()

	utils.WriteJSON(w, http.StatusOK, users)
}

func (h *Handler) handleGetUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("ID")

	user, _ := h.store.GetUserById(id)

	utils.WriteJSON(w, http.StatusOK, user)
}
