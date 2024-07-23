package user

import (
	"database/sql"
	"test/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUser(user types.User) error {
	_, err := s.db.Exec("INSERT INTO users (id, name, password) VALUES (?,?,?)", user.ID, user.Name, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByName(name string) (*types.User, error) {
	row := s.db.QueryRow("SELECT id, password FROM users WHERE name = ?", name)

	var user types.User
	err := row.Scan(&user.ID, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUserById(id string) (*types.GetUser, error) {
	row := s.db.QueryRow("SELECT id, name FROM users WHERE id = ?", id)

	var user types.GetUser
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetAuthById(id string) (*types.Auth, error) {
	row := s.db.QueryRow("SELECT id, name FROM users WHERE id = ?", id)

	var user types.Auth
	err := row.Scan(&user.ID, &user.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (s *Store) GetUsers() ([]types.GetUsers, error) {
	rows, err := s.db.Query("SELECT id, name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.GetUsers
	for rows.Next() {
		var u types.GetUsers
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
