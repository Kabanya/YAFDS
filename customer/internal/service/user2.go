package service

import (
	"customer/models"
)

// абстракция чуть выше чем репо

type UserLoader interface {
	Load(string) (models.User, error)
}

type user2 struct {
	repo UserLoader
}

func NewUser2(repo UserLoader) *user2 {
	return &user2{repo: repo}
}

func (s /*service*/ *user2) Load(walletAddress string) (models.User, error) {
	return s.repo.Load(walletAddress)
}
