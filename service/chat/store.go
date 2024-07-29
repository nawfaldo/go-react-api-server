package chat

import (
	"database/sql"
	"fmt"
	"test/types"

	"github.com/google/uuid"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetMessages(chatID string) (any, error) {
	rows, err := s.db.Query("SELECT user_id, message FROM chat_messages WHERE chat_id = $1", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	type UM struct {
		UserName string `json:"name"`
		UserID   string `json:"user_id"`
		Message  string `json:"message"`
	}

	var messages []UM
	for rows.Next() {
		var um UM

		if err := rows.Scan(&um.UserID, &um.Message); err != nil {
			return nil, err
		}

		row := s.db.QueryRow("SELECT name FROM users WHERE id = $1", um.UserID)

		if err := row.Scan(&um.UserName); err != nil {
			return nil, err
		}

		um.UserID = ""

		messages = append(messages, um)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *Store) GetChatUsers(userID string, chatID string) any {
	rows, err := s.db.Query("SELECT user_id FROM chat_or_server_users WHERE chat_id = $1 AND user_id != $2", chatID, userID)
	if err != nil {
		return nil
	}
	defer rows.Close()

	type U struct {
		UserID   string `json:"user_id"`
		UserName string `json:"name"`
	}

	var users []U
	for rows.Next() {
		var u U

		if err := rows.Scan(&u.UserID); err != nil {
			return nil
		}

		row := s.db.QueryRow("SELECT name FROM users WHERE id = $1", u.UserID)

		if err := row.Scan(&u.UserName); err != nil {
			return nil
		}

		u.UserID = ""

		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil
	}

	return users
}

func (s *Store) CreateMessage(message types.Message) error {
	_, err := s.db.Exec("INSERT INTO chat_messages (message, user_id, chat_id) VALUES ($1, $2, $3)", message.Message, message.UserID, message.ChatID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetOrCreateChat(user types.GetOrCreateChatPayload) (string, error) {
	var chatID string

	query := `
		SELECT c.id
		FROM chats c
		JOIN chat_or_server_users cu ON c.id = cu.chat_id
		WHERE cu.user_id IN ($1, $2)
		GROUP BY c.id
		HAVING COUNT(DISTINCT cu.user_id) = 2 AND COUNT(cu.user_id) = 2;
	`

	row := s.db.QueryRow(query, user.UserOne, user.UserTwo)
	err := row.Scan(&chatID)
	if err != nil {
		if err == sql.ErrNoRows {
			chatID = uuid.New().String()

			_, _ = s.db.Exec("INSERT INTO chats (id, name, server_id, server_role_id, server_chat_category_id) VALUES ($1, NULL, NULL, NULL, NULL)", chatID)

			_, _ = s.db.Exec("INSERT INTO chat_or_server_users (user_id, chat_id, server_id, server_role_id) VALUES ($1, $2, NULL, NULL)", user.UserOne, chatID)
			_, _ = s.db.Exec("INSERT INTO chat_or_server_users (user_id, chat_id, server_id, server_role_id) VALUES ($1, $2, NULL, NULL)", user.UserTwo, chatID)
		}
	}

	return chatID, nil
}

func (s *Store) GetChatsByUserID(userID string) any {
	rows, err := s.db.Query("SELECT chat_id FROM chat_or_server_users WHERE user_id = $1", userID)
	if err != nil {
		fmt.Println("err")
		return nil
	}
	defer rows.Close()

	type C struct {
		ChatID   string `json:"chat_id"`
		UserName string `json:"user_name"`
	}

	type U struct {
		ID   string `json:"user_id"`
		Name string `json:"name"`
	}

	var chats []C

	for rows.Next() {
		var c C
		if err := rows.Scan(&c.ChatID); err != nil {
			fmt.Println("err")
			return nil
		}

		var otherUser U

		row := s.db.QueryRow("SELECT user_id FROM chat_or_server_users WHERE chat_id = $1 AND user_id != $2", c.ChatID, userID)

		if err := row.Scan(&otherUser.ID); err != nil {
			fmt.Println("err")
			return nil
		}

		row = s.db.QueryRow("SELECT name FROM users WHERE id = $1", otherUser.ID)

		if err := row.Scan(&otherUser.Name); err != nil {
			fmt.Println("err")
			return nil
		}

		c.UserName = otherUser.Name

		chats = append(chats, c)
	}

	return chats
}
