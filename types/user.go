package types

import "time"

type UserStore interface {
	CreateUser(User) error
	GetUserByName(name string) (*User, error)
	GetUserById(id string) (*GetUser, error)
	GetAuthById(id string) (*Auth, error)
	GetUsers() ([]GetUsers, error)
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Auth struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetUsers struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RegisterUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginUserPayload struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type GetUserByIdPayload struct {
	ID string `json:"id" validate:"required"`
}
