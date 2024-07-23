package types

type ChatStore interface{}

type Chat struct {
	UserID  string `json:"user_id"`
	ChatID  string `json:"chat_id"`
	Message string `json:"message"`
}
