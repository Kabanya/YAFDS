package service

import (
	"github.com/google/uuid"
)

// абстракция чуть выше чем репо

type UserSaver interface {
	Save(uuid.UUID, string, string, string) error
}

type User interface {
	Save(uuid.UUID, string, string, string) error
	// Load(walletAddress string) (modules.User, error)
}

type user struct {
	repo UserSaver
}

func NewUser1(repo UserSaver) *user {
	return &user{repo: repo}
}

func (s /*service*/ *user) Save(id uuid.UUID, name string, walletAddress string, address string) error {

	return s.repo.Save(id, name, walletAddress, address)
}
