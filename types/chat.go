package types

type ChatStore interface {
	CreateMessage(Message) error
	GetMessages(chat_id string) (any, error)
	GetOrCreateChat(GetOrCreateChatPayload) (string, error)
	GetChatsByUserID(user_id string) any
	GetChatUsers(user_id string, chat_id string) any
}

type Chat struct {
	ID string `json:"id"`
}

type Message struct {
	UserID   string `json:"user_id"`
	UserName string `json:"name"`
	ChatID   string `json:"chat_id"`
	Message  string `json:"message"`
}

type SendMessagePayload struct {
	ChatID  string `json:"chat_id" validate:"required"`
	Message string `json:"message" validate:"required"`
	Name    string `json:"name" validate:"required"`
}

type GetOrCreateChatPayload struct {
	UserOne string `json:"UserOne" validate:"required"`
	UserTwo string `json:"UserTwo" validate:"required"`
}
